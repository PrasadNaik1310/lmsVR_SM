package handlers

import (
	"log"
	"net/http"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func DeleteModule(c *gin.Context) {

	log.Printf("Delete module request received. module_id=%s", c.Param("id"))

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

	var lessonCount int64

	if err := db.DB.
		Model(&models.Lesson{}).
		Where("module_id = ?", moduleID).
		Count(&lessonCount).Error; err != nil {

		log.Printf(
			"Failed counting lessons for module %s : %v",
			moduleID,
			err,
		)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to validate module",
		})
		return
	}

	log.Printf(
		"Module %s contains %d lessons",
		moduleID,
		lessonCount,
	)

	if lessonCount > 0 {

		log.Printf(
			"Delete denied. module_id=%s contains lessons",
			moduleID,
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot delete module containing lessons",
		})
		return
	}

	if err := db.DB.Delete(&module).Error; err != nil {

		log.Printf("Failed deleting module %s : %v", moduleID, err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete module",
		})
		return
	}

	log.Printf("Module deleted successfully. module_id=%s", moduleID)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
