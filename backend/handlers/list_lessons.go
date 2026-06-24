package handlers

import (
	"log"
	"net/http"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListLessons(c *gin.Context) {

	log.Printf(
		"List lessons request received. module_id=%s",
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

	var lessons []models.Lesson

	if err := db.DB.
		Where("module_id = ?", moduleID).
		Order("position ASC").
		Find(&lessons).Error; err != nil {

		log.Printf(
			"Failed fetching lessons for module %s : %v",
			moduleID,
			err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch lessons",
		})
		return
	}

	log.Printf(
		"Lessons fetched successfully. module_id=%s count=%d",
		moduleID,
		len(lessons),
	)

	c.JSON(http.StatusOK, gin.H{
		"lessons": lessons,
	})
}
