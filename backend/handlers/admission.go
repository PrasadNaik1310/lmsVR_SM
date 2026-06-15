package handlers

import (
	//"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/requests"
	"github.com/PrasadNaik1310/LMSVR_SM/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateEnquiry creates a new enquiry for an interested course
func CreateEnquiry(c *gin.Context) {
	var req requests.CreateEnquiryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse course ID
	courseID, err := uuid.Parse(req.InterestedCourseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid interested_course_id"})
		return
	}

	// Verify course exists
	var course models.Course
	if err := db.DB.Where("id = ?", courseID).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		log.Println("Course not exists ")
		return
	}

	enquiry := models.Enquiry{
		ID:                 uuid.New(),
		FullName:           req.FullName,
		Email:              req.Email,
		Phone:              req.Phone,
		InterestedCourseID: courseID,
		Status:             "new",
		Notes:              req.Notes,
		CreatedAt:          time.Now(),
	}

	if err := db.DB.Create(&enquiry).Error; err != nil {
		log.Printf("Failed to create enquiry: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create enquiry"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"enquiry": gin.H{
			"id":         enquiry.ID,
			"status":     enquiry.Status,
			"created_at": enquiry.CreatedAt,
		},
	})
}

// ListEnquiries lists all enquiries with pagination and optional filters
func ListEnquiries(c *gin.Context) {
	page := 1
	size := 20
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if s := c.Query("size"); s != "" {
		if parsed, err := strconv.Atoi(s); err == nil && parsed > 0 {
			size = parsed
		}
	}

	status := c.Query("status")
	courseID := c.Query("course_id")

	query := db.DB.Model(&models.Enquiry{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if courseID != "" {
		if cid, err := uuid.Parse(courseID); err == nil {
			query = query.Where("interested_course_id = ?", cid)
		}
	}

	var total int64
	query.Count(&total)

	var enquiries []models.Enquiry
	if err := query.Offset((page - 1) * size).Limit(size).Find(&enquiries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch enquiries"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enquiries": enquiries,
		"page":      page,
		"size":      size,
		"total":     total,
	})
}

// GetEnquiry retrieves a single enquiry by ID
func GetEnquiry(c *gin.Context) {
	enquiryIDStr := c.Param("id")
	enquiryID, err := uuid.Parse(enquiryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid enquiry id"})
		return
	}

	var enquiry models.Enquiry
	if err := db.DB.Where("id = ?", enquiryID).First(&enquiry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "enquiry not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"enquiry": enquiry})
}

// UpdateEnquiryStatus updates the status of an enquiry
func UpdateEnquiryStatus(c *gin.Context) {
	enquiryIDStr := c.Param("id")
	enquiryID, err := uuid.Parse(enquiryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid enquiry id"})
		return
	}

	var req requests.UpdateEnquiryStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate status
	allowedStatuses := map[string]bool{"new": true, "contacted": true, "approved": true, "rejected": true}
	if !allowedStatuses[req.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	var enquiry models.Enquiry
	if err := db.DB.Where("id = ?", enquiryID).First(&enquiry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "enquiry not found"})
		return
	}

	if err := db.DB.Model(&enquiry).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update enquiry status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"enquiry": gin.H{
			"id":     enquiry.ID,
			"status": req.Status,
		},
	})
}

// CreateApplication creates a new application for an enquiry
func CreateApplication(c *gin.Context) {
	var req requests.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse UUIDs
	enquiryID, err := uuid.Parse(req.EnquiryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid enquiry_id"})
		return
	}

	courseID, err := uuid.Parse(req.AppliedCourseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid applied_course_id"})
		return
	}

	// Verify enquiry exists
	var enquiry models.Enquiry
	if err := db.DB.Where("id = ?", enquiryID).First(&enquiry).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "enquiry not found"})
		return
	}

	// Verify course exists
	var course models.Course
	if err := db.DB.Where("id = ?", courseID).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}

	application := models.Application{
		ID:                uuid.New(),
		EnquiryID:         enquiryID,
		AppliedCourseID:   courseID,
		ApplicationStatus: "pending",
		SubmittedAt:       time.Now(),
	}

	if err := db.DB.Create(&application).Error; err != nil {
		log.Printf("Failed to create application: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create application"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"application": gin.H{
			"id":     application.ID,
			"status": application.ApplicationStatus,
		},
	})
}

// ListApplications lists all applications with pagination and optional filters
func ListApplications(c *gin.Context) {
	page := 1
	size := 20
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if s := c.Query("size"); s != "" {

		if parsed, err := strconv.Atoi(s); err == nil && parsed > 0 {
			if parsed > 100 {
				size = 100
			}
			size = parsed

		}
	}

	status := c.Query("status")
	courseID := c.Query("course_id")

	query := db.DB.Model(&models.Application{})
	if status != "" {
		query = query.Where("application_status = ?", status)
	}
	if courseID != "" {
		if cid, err := uuid.Parse(courseID); err == nil {
			query = query.Where("applied_course_id = ?", cid)
		}
	}

	var total int64
	query.Count(&total)

	var applications []models.Application
	if err := query.Offset((page - 1) * size).Limit(size).Find(&applications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch applications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"applications": applications,
		"page":         page,
		"size":         size,
		"total":        total,
	})
}

// GetApplication retrieves a single application by ID
func GetApplication(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		return
	}

	var application models.Application
	if err := db.DB.Where("id = ?", appID).First(&application).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		log.Println("Application Not found , Get application !!")
		return
	}

	c.JSON(http.StatusOK, gin.H{"application": application})
}

// generateTemporaryPassword generates a random temporary password
func generateTemporaryPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(password)
}

// ApproveApplication approves an application and creates a user/student account
func ApproveApplication(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		log.Printf("Application ID %s invalid :(", appID)
		return
	}

	var req requests.ApproveApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error , request not bind ": err.Error()})
		log.Printf("request format invalid , admission.go err dekh lee :-> %s", err)
		return
	}

	// Start transaction
	/*tx := db.DB.Begin(c,&sql.TxOptions{
		Isolation:sql.LevelSerializable,
		ReadOnly: false,
	})*/
	//tx := db.DB.BeginTx(c, nil)
	tx := db.DB.Begin()
	log.Printf("TX started for id %s", appID)
	//rollback

	// Fetch application
	var application models.Application
	if err := tx.Where("id = ?", appID).First(&application).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found , application nahi haiii ..."})
		log.Println("Application nahi milaaa bitch!!")
		return
	}

	if application.ApplicationStatus == "approved" {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "application is not in pending status"})
		log.Println("Application pending hai hi nahii bc")
		return
	}

	// Fetch enquiry
	var enquiry models.Enquiry
	if err := tx.Where("id = ?", application.EnquiryID).First(&enquiry).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "associated enquiry not found"})
		log.Println("Enquiry not found !!")

		return
	}

	/*// Generate temporary password if not provided
	tempPassword := req.TemporaryPassword
	if tempPassword == "" {
		tempPassword = generateTemporaryPassword(12)
	}*/
	hashedPassword, err := services.HashPassword(generateTemporaryPassword(12))
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to secure temporary password"})
		return
	}

	// Create user
	userID := uuid.New()
	// Parse full name into first and last name (split on first space)
	firstName := enquiry.FullName
	lastName := ""
	if idx := findFirstSpace(enquiry.FullName); idx != -1 {
		firstName = enquiry.FullName[:idx]
		lastName = enquiry.FullName[idx+1:]
	}

	// Determine student role (assuming student role exists, get it or use a default)
	var studentRole models.Role
	if err := tx.Where("name = ?", "student").First(&studentRole).Error; err != nil {
		// If student role doesn't exist, create a default one
		studentRole.ID = uuid.New()
		studentRole.Name = "student"
		studentRole.Description = "Student role"
		if err := tx.Create(&studentRole).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student role"})
			return
		}
	}

	user := models.User{
		ID:           userID,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        enquiry.Email,
		PasswordHash: hashedPassword,
		Phone:        enquiry.Phone,
		RoleID:       studentRole.ID,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to create user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user account"})
		return
	}
	// geenrating Enrollement number logic
	var enrollmentSeq int64

	if err := tx.Raw(
		"SELECT nextval('student_enrollment_seq')",
	).Scan(&enrollmentSeq).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate enrollment number",
		})
		return
	}
	studentID := uuid.New()

	// Create student profile
	student := models.Student{
		ID:               studentID,
		UserID:           userID,
		EnrollmentNumber: fmt.Sprintf("%d %d", time.Now().Year(), enrollmentSeq),
		AdmissionDate:    time.Now(),
		CreatedAt:        time.Now(),
	}

	if err := tx.Create(&student).Error; err != nil {
		tx.Rollback()
		log.Printf("Failed to create student profile: %v , Rollback TX id :%v", err, tx)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create student profile"})
		return
	}

	// Update application status
	if err := tx.Model(&application).Update("application_status", "approved").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update application status"})
		log.Printf("Failed to accept application %s", application)
		return
	}

	// Update enquiry status
	if err := tx.Model(&enquiry).Update("status", "approved").Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update enquiry status"})
		return
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to complete approval"})
		return
	}

	// TODO: Send welcome email (stubbed for now)
	log.Printf("Welcome email would be sent to %s with temporary password", enquiry.Email)
	log.Printf("Temp password %s", hashedPassword)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"application": gin.H{
			"id":     application.ID,
			"status": "approved",
			//"temporary_password": tempPassword,
		},
	})
}

// RejectApplication rejects an application
func RejectApplication(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid application id"})
		log.Printf("Application id %s not found ", appID)
		return
	}

	var req requests.RejectApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Printf("Application id %s request not bind . Err %s ", appID, err)
		return
	}

	var application models.Application
	if err := db.DB.Where("id = ?", appID).First(&application).Error; err != nil {
		log.Printf("Application not found !!")
		c.JSON(http.StatusNotFound, gin.H{"error": "application not found"})
		return
	}

	if application.ApplicationStatus == "rejected" || application.ApplicationStatus == "approved" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "application is not in pending status"})
		log.Printf("Application %s is not pending", appID)
		return
	}

	// Update application status and remarks
	if err := db.DB.Model(&application).Updates(map[string]interface{}{
		"application_status": "rejected",
		//"remarks":            req.Remarks,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reject application"})
		log.Printf("Failed to reject the application %s", application)
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"application": gin.H{
			"id":     application.ID,
			"status": "rejected",
		},
	})
}

// Helper function to find first space in a string
func findFirstSpace(s string) int {
	for i, r := range s {
		if r == ' ' {
			return i
		}
	}
	return -1
}

var now = time.Now()
