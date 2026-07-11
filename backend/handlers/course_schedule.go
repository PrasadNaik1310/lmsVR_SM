package handlers

import (
	"errors"
	"log"
	"net/http"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetSchedule(c *gin.Context) {
	log.Printf("list schedule request recieved with id %v", c.Param("id"))
	scheduleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Error : Invalid schedule ID ")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid schedule id"})
		return
	}
	var schedule models.CourseSchedule
	if err := db.DB.Where("id = ?", scheduleID).First(&schedule).Error; err != nil {
		//Trying a new error handling technique ,
		//using errors.Is travsers through stack trace and find the exact error specified in the function
		if errors.Is(err, gorm.ErrRecordNotFound) {

			log.Println("Error:Schedules not found in DB")
			c.JSON(http.StatusNotFound, gin.H{"error": "schedule not found"})
			return
		}
		log.Printf("Failed to fetch schedules :%v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch schedule"})

		return
	}
	c.JSON(http.StatusOK, gin.H{"schedule": schedule})
}
