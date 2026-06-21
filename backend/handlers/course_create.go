package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/requests"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateCourse(c *gin.Context) {

	var req requests.CreateCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error: Bad request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	sessionID, err := uuid.Parse(req.AcademicSessionID)
	if err != nil {
		log.Printf("Error: Invalid academic_session_id")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid academic_session_id",
		})
		return
	}

	var session models.AcademicSession
	if err := db.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		log.Printf("Error: academic session not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "academic session not found",
		})
		return
	}

	var startDate *time.Time
	var endDate *time.Time

	if req.StartDate != "" {
		t, err := time.Parse("2006-01-02", req.StartDate) //2006-01-02 go's special date format reference

		if err != nil {
			log.Printf("Error: Invalid Start_date")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid start_date",
			})
			return
		}
		startDate = &t
	}

	if req.EndDate != "" {
		t, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			log.Printf("Error: Invalid end_date")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid end_date",
			})
			return
		}
		endDate = &t
	}

	if endDate.Before(*startDate) || endDate.Equal(*startDate) {
		log.Printf("Error : endDates before or equal startDate")
		c.JSON(http.StatusBadRequest, gin.H{"error": "End Date must be after the start Date"})

		return
	}
	/*userID := c.GetString("user_id")

	creatorID, err := uuid.Parse(userID)*/

	CreatorIDAny, exists := c.Get("user_id")
	if !exists {
		log.Println("User ID not found in request ")
		c.JSON(http.StatusUnauthorized, gin.H{
			"Error": "User ID not Found ",
		})
	}
	creatorID, ok := CreatorIDAny.(uuid.UUID)
	if !ok {
		log.Printf("Invalid user ID type %v", creatorID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid User ID type ",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid user",
		})
		return
	}

	course := models.Course{
		ID:                uuid.New(),
		Title:             req.Title,
		Description:       req.Description,
		Level:             req.Level,
		Status:            "DRAFT",
		TotalSeats:        req.TotalSeats,
		BookedSeats:       0, //TODO: making 0 for testing till Enrolled module is created.
		StartDate:         startDate,
		EndDate:           endDate,
		MeetLink:          req.MeetLink,
		CreatedBy:         creatorID,
		AcademicSessionID: sessionID,
	}

	if err := db.DB.Create(&course).Error; err != nil {
		log.Println("Error:Failed to write course in DB .")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create course",
		})
		return
	}

	log.Println("Course Written on DB . ")
	c.JSON(http.StatusCreated, gin.H{
		"course_id": course.ID,
		"status":    course.Status,
	})
	log.Print(course)
}
