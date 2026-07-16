package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/requests"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateSchedule(c *gin.Context) {
	log.Printf("Request received to create schedule for course id = %v", c.Param("id"))
	var req requests.CreateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Request not binding with request struct :-> %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return

	}
	log.Printf("Create Schedule request recieved with course id = %s", c.Param("id"))
	CourseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Invalid Course ID : %s", c.Param("id"))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}
	var Course models.Course

	if err := db.DB.First(&Course, "id = ?", CourseID).Error; err != nil {
		log.Printf("Course not found. course_id=%s", CourseID)

		c.JSON(http.StatusNotFound, gin.H{
			"error": "course not found",
		})
		return
	}
	if Course.Status != "PUBLISHED" {
		log.Printf("Course not published ")
		c.JSON(http.StatusFailedDependency, gin.H{
			"error": "course is not published",
		})
		return
	}

	var lesson models.Lesson
	if err := db.DB.First(&lesson, "id = ?", req.LessonID).Error; err != nil {
		log.Printf("Lesson not found err:-> %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "lesson not found",
		})
		return
	}

	var module models.CourseModule
	if err := db.DB.First(&module, "id = ?", lesson.ModuleID).Error; err != nil {
		log.Printf("Module not found for lesson err:-> %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "lesson module not found",
		})
		return
	}

	if module.CourseID != CourseID {
		log.Printf("Lesson module does not belong to course. lesson_id=%s course_id=%s module_course_id=%s", req.LessonID, CourseID, module.CourseID)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "lesson belongs to another course",
		})
		return
	}

	var teacher models.Teacher
	log.Printf("TeacherID received from request = %s", req.TeacherID)
	if err := db.DB.First(&teacher, "id = ?", req.TeacherID).Error; err != nil {
		log.Printf("Teacher not found err:-> %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"error": "teacher not found",
		})
		return
	}

	plannedDate, err := time.Parse("2006-01-02", req.PlannedDate)
	if err != nil {
		log.Printf("Invalid planned_date: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid planned_date",
		})
		return
	}

	startTime, err := time.Parse("15:04", req.StartTime)
	if err != nil {
		log.Printf("Invalid start_time: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid start_time",
		})
		return
	}

	endTime, err := time.Parse("15:04", req.EndTime)
	if err != nil {
		log.Printf("Invalid end_time: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid end_time",
		})
		return
	}

	if !endTime.After(startTime) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "end time must be greater than start time",
		})
		return
	}

	var existingSchedule models.CourseSchedule
	if err := db.DB.Where(
		"lesson_id = ? AND status NOT IN ?",
		req.LessonID,
		[]string{"COMPLETED", "CANCELLED"},
	).First(&existingSchedule).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "lesson already has an active schedule",
		})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Failed to check existing schedule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to validate existing schedule",
		})
		return
	}

	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user id not found",
		})
		return
	}

	createdBy, ok := userIDAny.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id type",
		})
		return
	}

	schedule := models.CourseSchedule{
		ID:               uuid.New(),
		CourseID:         CourseID,
		LessonID:         req.LessonID,
		TeacherID:        req.TeacherID,
		PlannedDate:      plannedDate,
		PlannedStartTime: req.StartTime,
		PlannedEndTime:   req.EndTime,
		Status:           "SCHEDULED",
		CreatedBy:        createdBy,
	}

	if err := db.DB.Create(&schedule).Error; err != nil {
		log.Printf("Failed to create schedule: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create schedule",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"schedule_id": schedule.ID,
		"status":      schedule.Status,
	})

}
