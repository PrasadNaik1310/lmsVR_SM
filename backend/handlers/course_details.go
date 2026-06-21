package handlers

import (
	"log"
	"net/http"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetCourseDetails(c *gin.Context) {

	id := c.Param("id")

	courseID, err := uuid.Parse(id)
	if err != nil {
		log.Println("Error : Invalid course ID ")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	var course models.Course

	if err := db.DB.First(&course, "id = ?", courseID).Error; err != nil {
		log.Println("Error :Course Not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "course not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           course.ID,
		"title":        course.Title,
		"description":  course.Description,
		"level":        course.Level,
		"status":       course.Status,
		"start_date":   course.StartDate,
		"end_date":     course.EndDate,
		"total_seats":  course.TotalSeats,
		"booked_seats": course.BookedSeats,
		"meet_link":    course.MeetLink,
	})
	log.Print(course)
}
