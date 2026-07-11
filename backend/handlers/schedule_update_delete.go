package handlers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/requests"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UpdateSchedule(c *gin.Context) {
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

	var req requests.UpdateScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.TeacherID != "" {
		teacherID, parseErr := uuid.Parse(req.TeacherID)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid teacher_id"})
			return
		}

		var teacher models.Teacher
		if err := db.DB.First(&teacher, "id = ?", teacherID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "teacher not found"})
			return
		}
		schedule.TeacherID = teacherID
	}

	if req.PlannedDate != "" {
		plannedDate, parseErr := time.Parse("2006-01-02", req.PlannedDate)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid planned_date"})
			return
		}
		schedule.PlannedDate = plannedDate
	}

	startTime := schedule.PlannedStartTime
	endTime := schedule.PlannedEndTime
	if req.StartTime != "" {
		startTime = req.StartTime
	}
	if req.EndTime != "" {
		endTime = req.EndTime
	}

	parsedStart, err := time.Parse("15:04", startTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time"})
		return
	}
	parsedEnd, err := time.Parse("15:04", endTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time"})
		return
	}
	if !parsedEnd.After(parsedStart) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_time must be greater than start_time"})
		return
	}
	schedule.PlannedStartTime = startTime
	schedule.PlannedEndTime = endTime

	if req.Status != "" {
		newStatus := strings.ToUpper(strings.TrimSpace(req.Status))
		currentStatus := strings.ToUpper(strings.TrimSpace(schedule.Status))
		if !isValidScheduleStatus(newStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}
		if !isAllowedScheduleTransition(currentStatus, newStatus) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status transition"})
			return
		}
		schedule.Status = newStatus
	}

	if err := db.DB.Save(&schedule).Error; err != nil {
		log.Printf("Failed to update schedule %s: %v", scheduleID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedule": schedule})
}

func DeleteSchedule(c *gin.Context) {
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

	if strings.EqualFold(schedule.Status, "COMPLETED") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "completed schedules cannot be deleted"})
		return
	}

	if err := db.DB.Delete(&schedule).Error; err != nil {
		log.Printf("Failed to delete schedule %s: %v", scheduleID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func isValidScheduleStatus(status string) bool {
	validStatuses := map[string]struct{}{
		"SCHEDULED": {},
		"ONGOING":   {},
		"COMPLETED": {},
		"POSTPONED": {},
		"CANCELLED": {},
	}
	_, ok := validStatuses[status]
	return ok
}

func isAllowedScheduleTransition(from, to string) bool {
	if from == to {
		return true
	}

	allowed := map[string]map[string]struct{}{
		"SCHEDULED": {
			"ONGOING":   {},
			"POSTPONED": {},
			"CANCELLED": {},
		},
		"ONGOING": {
			"COMPLETED": {},
			"POSTPONED": {},
			"CANCELLED": {},
		},
		"POSTPONED": {
			"SCHEDULED": {},
			"CANCELLED": {},
		},
	}

	next, ok := allowed[from]
	if !ok {
		return false
	}

	_, ok = next[to]
	return ok
}
