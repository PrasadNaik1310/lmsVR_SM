package handlers

import (
	"log"
	"net/http"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeleteLesson(c *gin.Context) {

	log.Printf(
		"Delete lesson request received. lesson_id=%s",
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

	if err := db.DB.Delete(&lesson).Error; err != nil {

		log.Printf(
			"Failed deleting lesson %s : %v",
			lessonID,
			err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete lesson",
		})
		return
	}

	log.Printf(
		"Lesson deleted successfully. lesson_id=%s",
		lessonID,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
