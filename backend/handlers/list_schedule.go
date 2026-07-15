package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/requests"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListSchedules(c *gin.Context) {
	log.Printf("List schedules request recieved with course id = %v", c.Param("id"))
	CourseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Failed to Parse UUID , invalid UUID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid uuid, failed to parse",
		})
		return
	}

	page := 1
	size := 20
	if p := c.Query("page"); p != "" {
		if parsed, parseErr := strconv.Atoi(p); parseErr == nil && parsed > 0 {
			page = parsed
		}
	}
	if s := c.Query("size"); s != "" {
		if parsed, parseErr := strconv.Atoi(s); parseErr == nil && parsed > 0 {
			size = parsed
		}
	}

	status := c.Query("status")
	plannedDate := c.Query("planned_date")
	teacherID := c.Query("teacher_id")

	query := db.DB.Model(&models.CourseSchedule{}).Where("course_id = ?", CourseID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if plannedDate != "" {
		if parsedDate, parseErr := time.Parse("2006-01-02", plannedDate); parseErr == nil {
			query = query.Where("planned_date::date = ?", parsedDate.Format("2006-01-02"))
		}
	}
	if teacherID != "" {
		if parsedTeacherID, parseErr := uuid.Parse(teacherID); parseErr == nil {
			query = query.Where("teacher_id = ?", parsedTeacherID)
		}
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		log.Printf("Error: Failed to count schedules %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to count schedules",
		})
		return
	}

	var schedule []models.CourseSchedule
	if err := query.Order("planned_date DESC").Offset((page - 1) * size).Limit(size).Find(&schedule).Error; err != nil {
		log.Printf("Error: Failed to Fetch Schedules  %v ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch schedules",
		})
		return

	}
	var response []requests.CourseScheduleResponse

	for _, schedule := range schedules{
		response = append(response, ...)
	}
	log.Printf(
		"Schedules fetched successfully. course_id=%s count=%d",
		CourseID,
		len(schedule),
	)

	c.JSON(http.StatusOK, gin.H{
		"schedules": response,
		"page":      page,
		"size":      size,
		"total":     total,
	})
}
