package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/requests"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateSessionLog(c *gin.Context) {
	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	var schedule models.CourseSchedule
	if err := db.DB.First(&schedule, "id = ?", scheduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "schedule not found"})
		return
	}

	var existing models.CourseLog
	if err := db.DB.First(&existing, "schedule_id = ?", scheduleID).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "session log already exists for this schedule"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate session log"})
		return
	}

	var req requests.CreateCourseLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	completionStatus := strings.ToUpper(strings.TrimSpace(req.CompletionStatus))
	if !isValidLogCompletionStatus(completionStatus) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid completion_status"})
		return
	}

	conductedDate, err := time.Parse("2006-01-02", req.ConductedDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conducted_date"})
		return
	}

	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id not found in context"})
		return
	}
	recordedBy, ok := userIDAny.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id type"})
		return
	}

	courseLog := models.CourseLog{
		ID:               uuid.New(),
		ScheduleID:       scheduleID,
		ConductedDate:    conductedDate,
		CompletionStatus: completionStatus,
		Remarks:          req.Remarks,
		Homework:         req.Homework,
		NextTopic:        req.NextTopic,
		RecordedBy:       recordedBy,
	}

	if err := db.DB.Create(&courseLog).Error; err != nil {
		log.Printf("Failed to create course log for schedule %s: %v", scheduleID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create session log"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"log_id": courseLog.ID})
}

func GetSessionLog(c *gin.Context) {
	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}

	var schedule models.CourseSchedule
	if err := db.DB.First(&schedule, "id = ?", scheduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "schedule not found"})
		return
	}

	var courseLog models.CourseLog
	if err := db.DB.First(&courseLog, "schedule_id = ?", scheduleID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "session log not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch session log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"log": courseLog})
}

func UpdateSessionLog(c *gin.Context) {
	logID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid log id"})
		return
	}

	var courseLog models.CourseLog
	if err := db.DB.First(&courseLog, "id = ?", logID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session log not found"})
		return
	}

	var req requests.UpdateCourseLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.CompletionStatus != "" {
		completionStatus := strings.ToUpper(strings.TrimSpace(req.CompletionStatus))
		if !isValidLogCompletionStatus(completionStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid completion_status"})
			return
		}
		courseLog.CompletionStatus = completionStatus
	}

	if req.Remarks != "" {
		courseLog.Remarks = req.Remarks
	}
	if req.Homework != "" {
		courseLog.Homework = req.Homework
	}
	if req.NextTopic != "" {
		courseLog.NextTopic = req.NextTopic
	}

	if err := db.DB.Save(&courseLog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update session log"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"log": courseLog})
}

func isValidLogCompletionStatus(status string) bool {
	validStatuses := map[string]struct{}{
		"COMPLETED":           {},
		"PARTIALLY_COMPLETED": {},
		"CANCELLED":           {},
	}

	_, ok := validStatuses[status]
	return ok
}
