package models

import (
	"time"

	"github.com/google/uuid"
)

type Course struct {
	ID                uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	Title             string     `gorm:"size:255;not null" json:"title"`
	Description       string     `gorm:"type:text" json:"description"`
	Level             string     `gorm:"size:50" json:"level"`
	ThumbnailURL      string     `gorm:"type:text" json:"thumbnail_url"`
	Status            string     `gorm:"size:50" json:"status"`
	TotalSeats        int        `gorm:"check:total_seats >= 0" json:"total_seats"`
	BookedSeats       int        `gorm:"default:0" json:"booked_seats"`
	StartDate         *time.Time `json:"start_date,omitempty"`
	EndDate           *time.Time `json:"end_date,omitempty"`
	MeetLink          string     `gorm:"type:text" json:"meet_link"`
	CreatedBy         uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	AcademicSessionID uuid.UUID  `gorm:"type:uuid" json:"academic_session_id"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type Batch struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CourseID    uuid.UUID  `gorm:"type:uuid;unique" json:"course_id"`
	BatchName   string     `gorm:"size:100" json:"batch_name"`
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	MaxStudents int        `gorm:"check:max_students > 0" json:"max_students"`
	Status      string     `gorm:"size:50" json:"status"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

type BatchTeacher struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	BatchID   uuid.UUID `gorm:"type:uuid" json:"batch_id"`
	TeacherID uuid.UUID `gorm:"type:uuid" json:"teacher_id"`
}
