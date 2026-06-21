package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GenerateCourseInvite(c *gin.Context) {

	log.Printf("Generate invite request received. course_id=%s", c.Param("id"))

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

	userIDAny, exists := c.Get("user_id")
	if !exists {
		log.Println("User not found in context")

		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	userID := userIDAny.(uuid.UUID)

	inviteCode := uuid.New().String()[:8]

	invite := models.CourseInvite{
		ID:         uuid.New(),
		CourseID:   courseID,
		InviteCode: inviteCode,
		CreatedBy:  userID,
		ExpiresAt:  nil,
	}

	if err := db.DB.Create(&invite).Error; err != nil {
		log.Printf("Failed creating invite. course_id=%s err=%v", courseID, err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate invite",
		})
		return
	}

	log.Printf("Invite generated successfully. invite_code=%s course_id=%s",
		inviteCode,
		courseID,
	)

	c.JSON(http.StatusCreated, gin.H{
		"invite_code":  inviteCode,
		"invite_link":  "http://localhost:5173/join/" + inviteCode, //TODO:Replace localhost link with the third party link
		"generated_at": time.Now().UTC(),
	})
}
