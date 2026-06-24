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

func CreateLesson(c *gin.Context) {

	log.Printf(
		"Create lesson request received. id=%s",
		c.Param("id"),
	)

	moduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {

		log.Printf("Invalid module id: %s", c.Param("id"))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid module id",
		})
		return
	}

	var module models.CourseModule

	if err := db.DB.First(&module, "id = ?", moduleID).Error; err != nil {

		log.Printf("Module not found. id=%s", moduleID)

		c.JSON(http.StatusNotFound, gin.H{
			"error": "module not found",
		})
		return
	}

	var req requests.CreateLessonRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		log.Printf("Invalid create lesson payload: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var lessonCount int64

	db.DB.Model(&models.Lesson{}).
		Where("id = ?", moduleID).
		Count(&lessonCount)

	lesson := models.Lesson{
		ID:           uuid.New(),
		ModuleID:     moduleID,
		Title:        req.Title,
		Description:  req.Description,
		ContentType:  req.ContentType,
		DurationMins: req.DurationMinutes,
		Position:     int(lessonCount) + 1,
		IsPublished:  false,
	}

	if err := db.DB.Create(&lesson).Error; err != nil {

		log.Printf("Failed creating lesson: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create lesson",
		})
		return
	}

	log.Printf(
		"Lesson created successfully. lesson_id=%s id=%s",
		lesson.ID,
		moduleID,
	)

	c.JSON(http.StatusCreated, gin.H{
		"lesson_id": lesson.ID,
	})
}
