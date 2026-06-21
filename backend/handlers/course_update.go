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

func UpdateCourse(c *gin.Context) {

	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Error: Invalid course ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	var course models.Course

	if err := db.DB.First(&course, "id = ?", courseID).Error; err != nil {
		log.Printf("Error: Course not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "course not found",
		})
		return
	}

	var req requests.UpdateCourseRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Error: Bad request")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Title != "" {
		course.Title = req.Title
	}

	if req.Description != "" {
		course.Description = req.Description
	}

	if req.Level != "" {
		course.Level = req.Level
	}

	if req.MeetLink != "" {
		course.MeetLink = req.MeetLink
	}

	if req.TotalSeats > 0 {
		course.TotalSeats = req.TotalSeats
	}

	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			log.Printf("Error: Invalid start Date ")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid start_date",
			})
			return
		}
		course.StartDate = &start
	}

	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			log.Printf("Error: Invalid End Date")
			c.JSON(http.StatusBadRequest, gin.H{

				"error": "invalid end_date",
			})
			return
		}
		course.EndDate = &end
	}

	if err := db.DB.Save(&course).Error; err != nil {
		log.Printf("Error: Failed to update course %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update course",
		})
		return
	}
	log.Printf("Request success ")
	log.Println(course)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
