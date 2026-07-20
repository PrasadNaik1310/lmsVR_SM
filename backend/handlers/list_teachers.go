package handlers

import (
	"net/http"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TeacherResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func ListTeachers(c *gin.Context) {
	var teachers []TeacherResponse

	err := db.DB.
		Table("teachers t").
		Select(`
			t.id,
			CONCAT(u.first_name, ' ', u.last_name) as name
		`).
		Joins("JOIN users u ON u.id = t.user_id").
		Scan(&teachers).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch teachers",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"teachers": teachers,
	})
}
