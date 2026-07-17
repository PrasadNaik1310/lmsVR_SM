package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
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

	query := db.DB.
		Table("course_schedules cs").
		Joins("JOIN lessons l ON l.id = cs.lesson_id").
		Joins("JOIN teachers t ON t.id = cs.teacher_id").
		Joins("JOIN users u ON u.id = t.user_id").
		Where("cs.course_id = ?", CourseID)

	if status != "" {
		query = query.Where("cs.status = ?", status)
	}

	if plannedDate != "" {
		if parsedDate, parseErr := time.Parse("2006-01-02", plannedDate); parseErr == nil {
			query = query.Where(
				"DATE(cs.planned_date) = ?",
				parsedDate.Format("2006-01-02"),
			)
		}
	}

	if teacherID != "" {
		if parsedTeacherID, parseErr := uuid.Parse(teacherID); parseErr == nil {
			query = query.Where("cs.teacher_id = ?", parsedTeacherID)
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

	var response []requests.CourseScheduleResponse

	if err := query.
		Select(`
			cs.id,
			cs.lesson_id,
			l.title AS lesson_title,
			cs.teacher_id,
			CONCAT(u.first_name, ' ', u.last_name) AS teacher_name,
			cs.planned_date,
			cs.planned_start_time,
			cs.planned_end_time,
			cs.status
		`).
		Order("cs.planned_date DESC").
		Offset((page - 1) * size).
		Limit(size).
		Scan(&response).Error; err != nil {

		log.Printf("Error: Failed to Fetch Schedules %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch schedules",
		})
		return
	}

	log.Printf(
		"Schedules fetched successfully. course_id=%s count=%d",
		CourseID,
		len(response),
	)

	c.JSON(http.StatusOK, gin.H{
		"schedules": response,
		"page":      page,
		"size":      size,
		"total":     total,
	})
}
