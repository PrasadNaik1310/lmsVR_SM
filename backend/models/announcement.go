package models

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Title      string    `gorm:"size:255;not null" json:"title"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	TargetType string    `gorm:"size:50" json:"target_type"`
	CreatedBy  uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type AnnouncementCourse struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	AnnouncementID uuid.UUID `gorm:"type:uuid" json:"announcement_id"`
	CourseID       uuid.UUID `gorm:"type:uuid" json:"course_id"`
}
