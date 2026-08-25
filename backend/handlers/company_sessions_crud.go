package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/PrasadNaik1310/LMSVR_SM/requests"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateAcademicSession(c *gin.Context) {
	var req requests.CreateAcademicSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
		return
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
		return
	}

	if !endDate.After(startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be after start_date"})
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	session := models.AcademicSession{
		ID:        uuid.New(),
		Name:      req.Name,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  isActive,
	}

	if err := db.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create academic session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"academic_session": session})
}

func ListAcademicSessions(c *gin.Context) {
	var sessions []models.AcademicSession
	if err := db.DB.Order("start_date asc").Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch academic sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"academic_sessions": sessions})
}

func GetAcademicSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}

	var session models.AcademicSession
	if err := db.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "academic session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch academic session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"academic_session": session})
}

func UpdateAcademicSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}

	var session models.AcademicSession
	if err := db.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "academic session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch academic session"})
		return
	}

	var req requests.UpdateAcademicSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Name != "" {
		session.Name = req.Name
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date"})
			return
		}
		session.StartDate = startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date"})
			return
		}
		session.EndDate = endDate
	}

	if !session.EndDate.After(session.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be after start_date"})
		return
	}

	if req.IsActive != nil {
		session.IsActive = *req.IsActive
	}

	if err := db.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update academic session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"academic_session": session})
}

func DeleteAcademicSession(c *gin.Context) {
	sessionID, err := uuid.Parse(c.Param("session_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}

	var session models.AcademicSession
	if err := db.DB.First(&session, "id = ?", sessionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "academic session not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch academic session"})
		return
	}

	if err := db.DB.Delete(&session).Error; err != nil {
		log.Printf("failed to delete academic session %s: %v", sessionID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete academic session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}
