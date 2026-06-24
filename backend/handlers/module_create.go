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

func CreateModule(c *gin.Context) {

	log.Printf("Create module request received. course_id=%s",
		c.Param("id"))

	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {

		log.Printf("Invalid course id: %s", c.Param("id"))

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	var course models.Course

	if err := db.DB.First(&course, "id = ?", courseID).Error; err != nil {

		log.Printf("Course not found. course_id=%s", courseID)

		c.JSON(http.StatusNotFound, gin.H{
			"error": "course not found",
		})
		return
	}

	var req requests.CreateModuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		log.Printf("Invalid create module payload: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	module := models.CourseModule{
		ID:          uuid.New(),
		CourseID:    courseID,
		Title:       req.Title,
		Description: req.Description,
		Position:    req.Position,
	}

	if err := db.DB.Create(&module).Error; err != nil {

		log.Printf("Failed creating module: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create module",
		})
		return
	}

	log.Printf(
		"Module created successfully. module_id=%s course_id=%s",
		module.ID,
		courseID,
	)

	c.JSON(http.StatusCreated, gin.H{
		"module_id": module.ID,
	})
}
