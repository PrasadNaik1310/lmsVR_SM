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

func UpdateModule(c *gin.Context) {

	log.Printf("Update module request received. module_id=%s", c.Param("id"))

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

		log.Printf("Module not found. module_id=%s", moduleID)

		c.JSON(http.StatusNotFound, gin.H{
			"error": "module not found",
		})
		return
	}

	var req requests.UpdateModuleRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		log.Printf("Invalid update module payload: %v", err)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if req.Title != "" {
		module.Title = req.Title
	}

	if req.Description != "" {
		module.Description = req.Description
	}

	if req.Position > 0 {
		module.Position = req.Position
	}

	if err := db.DB.Save(&module).Error; err != nil {

		log.Printf("Failed updating module %s : %v", moduleID, err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update module",
		})
		return
	}

	log.Printf(
		"Module updated successfully. module_id=%s title=%s",
		module.ID,
		module.Title,
	)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
