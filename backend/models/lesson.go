package models

import (
	"time"

	"github.com/google/uuid"
)

type CourseModule struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CourseID    uuid.UUID `gorm:"type:uuid" json:"course_id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Position    int       `gorm:"not null" json:"position"`
}

type Lesson struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ModuleID     uuid.UUID `gorm:"type:uuid" json:"module_id"`
	Title        string    `gorm:"size:255;not null" json:"title"`
	Description  string    `gorm:"type:text" json:"description"`
	ContentType  string    `gorm:"size:50" json:"content_type"`
	ContentURL   string    `gorm:"type:text" json:"content_url"`
	DurationMins int       `json:"duration_minutes"`
	Position     int       `json:"position"`
	IsPublished  bool      `gorm:"default:false" json:"is_published"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type LessonFile struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	LessonID   uuid.UUID `gorm:"type:uuid" json:"lesson_id"`
	FileName   string    `gorm:"size:255" json:"file_name"`
	FileURL    string    `gorm:"type:text;not null" json:"file_url"`
	UploadedBy uuid.UUID `gorm:"type:uuid" json:"uploaded_by"`
	UploadedAt time.Time `gorm:"autoCreateTime" json:"uploaded_at"`
}
