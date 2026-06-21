package handlers

import (
	"log"
	"net/http"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PublishCourse(c *gin.Context) {

	courseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Printf("Error: Invalid course ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid course id",
		})
		return
	}

	var course models.Course

	if err := db.DB.First(&course, "id = ?", courseID).Error; err != nil {
		log.Printf("Error: Course Not found")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "course not found",
		})
		return
	}

	var moduleCount int64

	db.DB.Model(&models.CourseModule{}).Where("course_id = ?", courseID).Count(&moduleCount)

	if moduleCount == 0 {
		log.Printf("Error: Course must contain at least one module")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "course must contain at least one module",
		})
		return
	}

	var lessonCount int64

	db.DB.Model(&models.Lesson{}).
		Joins("JOIN course_modules ON lessons.module_id = course_modules.id").
		Where("course_modules.course_id = ?", courseID).
		Count(&lessonCount)

	if lessonCount == 0 {
		log.Printf("Error: Course must contain at leat one lesson")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "course must contain at least one lesson",
		})
		return
	}

	course.Status = "PUBLISHED"

	if err := db.DB.Save(&course).Error; err != nil {
		log.Printf("Error: Failed to publish course , not able to save in DB")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to publish course",
		})
		return
	}
	log.Println(course)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"status":  course.Status,
	})
}
