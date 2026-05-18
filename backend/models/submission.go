package models

import (
	"time"

	"github.com/google/uuid"
)

type LessonSubmission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	LessonID    uuid.UUID `gorm:"type:uuid" json:"lesson_id"`
	StudentID   uuid.UUID `gorm:"type:uuid" json:"student_id"`
	FileURL     string    `gorm:"type:text;not null" json:"file_url"`
	SubmittedAt time.Time `gorm:"autoCreateTime" json:"submitted_at"`
	Remarks     string    `gorm:"type:text" json:"remarks"`
}
