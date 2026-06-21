package handlers

import (
	"log"
	"net/http"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetLessonDetails(c *gin.Context) {

	log.Printf(
		"Lesson details request received. lesson_id=%s",
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

	log.Printf(
		"Lesson fetched successfully. lesson_id=%s title=%s",
		lesson.ID,
		lesson.Title,
	)

	c.JSON(http.StatusOK, gin.H{
		"lesson": lesson,
	})
}
