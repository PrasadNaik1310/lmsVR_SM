package handlers

import (
	"log"
	"net/http"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/requests"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UpdateLesson(c *gin.Context) {

	log.Printf(
		"Update lesson request received. lesson_id=%s",
		c.Param("id"),
	)

	lessonID, err := uuid.Parse(c.Param("id"))
	if err != nil {

		log.Printf("Invalid lesson id: %s", c.Param("id"))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid lesson id",
		})
		return
	}

	var lesson models.Lesson

	if err := db.DB.
		First(&lesson, "id = ?", lessonID).Error; err != nil {

		log.Printf("Lesson not found. lesson_id=%s", lessonID)

		c.JSON(http.StatusNotFound, gin.H{
			"error": "lesson not found",
		})
		return
	}

	var req requests.UpdateLessonRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		log.Printf("Invalid lesson update payload: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Title != "" {
		lesson.Title = req.Title
	}

	if req.Description != "" {
		lesson.Description = req.Description
	}

	if req.ContentType != "" {
		lesson.ContentType = req.ContentType
	}

	if req.DurationMinutes > 0 {
		lesson.DurationMins = req.DurationMinutes
	}

	lesson.IsPublished = req.IsPublished

	if err := db.DB.Save(&lesson).Error; err != nil {

		log.Printf(
			"Failed updating lesson %s : %v",
			lessonID,
			err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update lesson",
		})
		return
	}

	log.Printf(
		"Lesson updated successfully. lesson_id=%s",
		lessonID,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
