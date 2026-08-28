package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/requests"
	"github.com/PrasadNaik1310/LMSVR_SM/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateCompanyUser(c *gin.Context) {
	var req requests.CreateCompanyUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	req.Role = strings.ToLower(strings.TrimSpace(req.Role))

	// Only roles that a company admin is allowed to create.
	allowedRoles := map[string]bool{
		"student": true,
		"teacher": true,
		"admin":   true,
	}

	if !allowedRoles[req.Role] {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid role; allowed roles are student, teacher, admin",
		})
		return
	}

	// Admin accounts are represented by the company_admin role.
	roleName := req.Role

	if req.Role == "admin" {
		roleName = "company_admin"
	}

	// Hash password BEFORE starting the DB transaction.
	hashedPassword, err := services.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to secure password",
		})
		return
	}

	tx := db.DB.Begin()

	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to start transaction",
		})
		return
	}

	// Rollback automatically if we return because of an error.
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	// Find role.
	var role models.Role

	if err := tx.
		Where("name = ?", roleName).
		First(&role).Error; err != nil {

		tx.Rollback()

		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "requested role does not exist",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to find role",
		})
		return
	}

	// Check duplicate email.
	var existingUser models.User

	if err := tx.
		Where("email = ?", req.Email).
		First(&existingUser).Error; err == nil {

		tx.Rollback()

		c.JSON(http.StatusConflict, gin.H{
			"error": "user with this email already exists",
		})
		log.Printf("Duplicate email")
		return
	} else if err != gorm.ErrRecordNotFound {

		tx.Rollback()

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to check existing user",
		})
		log.Printf("Teacher not found")
		return
	}

	userID := uuid.New()

	user := models.User{
		ID:           userID,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Phone:        req.Phone,
		RoleID:       role.ID,
		IsActive:     true,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		log.Printf("DB error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create user",
		})

		return
	}

	// Student-specific profile.
	if req.Role == "student" {

		if req.AdmissionDate == "" {
			tx.Rollback()

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "admission_date is required for students",
			})
			return
		}

		admissionDate, err := time.Parse("2006-01-02", req.AdmissionDate)

		if err != nil {
			tx.Rollback()

			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid admission_date; expected YYYY-MM-DD",
			})
			return
		}

		var dateOfBirth *time.Time

		if req.DateOfBirth != "" {
			parsedDOB, err := time.Parse(
				"2006-01-02",
				req.DateOfBirth,
			)

			if err != nil {
				tx.Rollback()

				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid date_of_birth; expected YYYY-MM-DD",
				})
				return
			}

			dateOfBirth = &parsedDOB
		}

		student := models.Student{
			ID:               uuid.New(),
			UserID:           userID,
			EnrollmentNumber: req.EnrollmentNumber,
			DateOfBirth:      dateOfBirth,
			Address:          req.Address,
			AdmissionDate:    admissionDate,
		}

		if err := tx.Create(&student).Error; err != nil {
			tx.Rollback()

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create student profile",
			})
			return
		}
	}

	// Teacher-specific profile.
	if req.Role == "teacher" {

		var joiningDate *time.Time

		if req.JoiningDate != "" {
			parsedJoiningDate, err := time.Parse(
				"2006-01-02",
				req.JoiningDate,
			)

			if err != nil {
				tx.Rollback()

				c.JSON(http.StatusBadRequest, gin.H{
					"error": "invalid joining_date; expected YYYY-MM-DD",
				})
				return
			}

			joiningDate = &parsedJoiningDate
		}

		teacher := models.Teacher{
			ID:             uuid.New(),
			UserID:         userID,
			Specialization: req.Specialization,
			Bio:            req.Bio,
			JoiningDate:    joiningDate,
		}

		if err := tx.Create(&teacher).Error; err != nil {
			tx.Rollback()

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to create teacher profile",
			})
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to commit user creation",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"user": gin.H{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"phone":      user.Phone,
			"role":       role.Name,
			"is_active":  user.IsActive,
		},
	})
}
func ListCompanyUsers(c *gin.Context) {
	type UserResponse struct {
		ID        uuid.UUID `json:"id"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name"`
		Email     string    `json:"email"`
		Phone     string    `json:"phone"`
		Role      string    `json:"role"`
		IsActive  bool      `json:"is_active"`
	}

	var users []UserResponse

	err := db.DB.
		Table("users").
		Select(`
			users.id,
			users.first_name,
			users.last_name,
			users.email,
			users.phone,
			roles.name AS role,
			users.is_active
		`).
		Joins("LEFT JOIN roles ON roles.id = users.role_id").
		Scan(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
