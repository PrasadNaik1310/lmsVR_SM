package handlers

import (
	"log"
	"net/http"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListModules(c *gin.Context) {

	log.Printf("List modules request received. with course 'id'=%s",
		c.Param("id"))

	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Error: Failed to parse Course ID , might be invalid ")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	var modules []models.CourseModule

	if err := db.DB.
		Where("course_id = ?", courseID).
		Order("position asc").
		Find(&modules).Error; err != nil {

		log.Printf("Failed fetching modules: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch modules",
		})
		return
	}

	log.Printf(
		"Modules fetched successfully. course_id=%s count=%d",
		courseID,
		len(modules),
	)

	c.JSON(http.StatusOK, gin.H{
		"modules": modules,
	})
}
