package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/PrasadNaik1310/LMSVR_SM/db"
	"github.com/PrasadNaik1310/LMSVR_SM/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AssignCourseToSession(c *gin.Context) {
	sessionIDStr := c.Param("session_id")
	courseIDStr := c.Param("course_id")

	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course_id"})
		return
	}

	var session models.AcademicSession
	if err := db.DB.Where("id = ?", sessionID).First(&session).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "academic session not found"})
		return
	}
	if !session.IsActive {
		c.JSON(http.StatusConflict, gin.H{"error": "target academic session is not active"})
		return
	}

	var course models.Course
	if err := db.DB.Where("id = ?", courseID).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}

	course.AcademicSessionID = sessionID
	if err := db.DB.Save(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign course to session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":             true,
		"course_id":           course.ID,
		"academic_session_id": session.ID,
		"assigned_at":         time.Now().UTC(),
	})
}

func ListCoursesBySession(c *gin.Context) {
	sessionIDStr := c.Param("session_id")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session_id"})
		return
	}
	//pagination
	page := 1
	size := 20
	if p := c.Query("page"); p != "" {
		// ignore errors, keep defaults on failure
		fmtS := 0
		fmtS, _ = strconv.Atoi(p)
		if fmtS > 0 {
			page = fmtS
		}
	}
	if s := c.Query("size"); s != "" {
		ss := 0
		ss, _ = strconv.Atoi(s)
		if ss > 0 {
			size = ss
		}
	}
	status := c.Query("status")

	var courses []models.Course
	query := db.DB.Model(&models.Course{}).Where("academic_session_id = ?", sessionID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	var total int64
	query.Count(&total)
	query = query.Offset((page - 1) * size).Limit(size)
	if err := query.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch courses"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"courses": courses, "page": page, "size": size, "total": total})
}

type createBatchRequest struct {
	BatchName   string `json:"batch_name" binding:"required"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date" binding:"required"`
	MaxStudents int    `json:"max_students" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

func CreateBatchForCourse(c *gin.Context) {
	courseIDStr := c.Param("course_id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course_id"})
		return
	}

	var req createBatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// parse dates
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format, use YYYY-MM-DD"})
		return
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format, use YYYY-MM-DD"})
		return
	}
	if end.Before(start) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be before or equal to end_date"})
		return
	}
	if req.MaxStudents <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "max_students must be greater than 0"})
		return
	}

	// check course exists
	var course models.Course
	if err := db.DB.Where("id = ?", courseID).First(&course).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
		return
	}

	// uniqueness: one batch per course
	var existing models.Batch
	if err := db.DB.Where("course_id = ?", courseID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "a batch already exists for this course"})
		return
	}

	batch := models.Batch{
		ID:          uuid.New(),
		CourseID:    courseID,
		BatchName:   req.BatchName,
		StartDate:   &start,
		EndDate:     &end,
		MaxStudents: req.MaxStudents,
		Status:      req.Status,
	}
	if err := db.DB.Create(&batch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create batch"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"batch": batch})
}

func GetBatchDetails(c *gin.Context) {
	courseIDStr := c.Param("course_id")
	batchIDStr := c.Param("batch_id")
	courseID, err := uuid.Parse(courseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid course_id"})
		return
	}
	batchID, err := uuid.Parse(batchIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch_id"})
		return
	}
	var batch models.Batch
	if err := db.DB.Where("id = ? AND course_id = ?", batchID, courseID).First(&batch).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "batch not found for given course"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"batch": batch})
}
