package models

import (
	"time"

	"github.com/google/uuid"
)

type LiveSession struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CourseID      uuid.UUID `gorm:"type:uuid" json:"course_id"`
	Title         string    `gorm:"size:255;not null" json:"title"`
	Description   string    `gorm:"type:text" json:"description"`
	SessionDate   time.Time `gorm:"not null" json:"session_date"`
	StartTime     string    `gorm:"type:time;not null" json:"start_time"`
	EndTime       string    `gorm:"type:time;not null" json:"end_time"`
	MeetLink      string    `gorm:"type:text" json:"meet_link"`
	SessionStatus string    `gorm:"size:50" json:"session_status"`
	CreatedBy     uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type LiveSessionAttendance struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	LiveSessionID    uuid.UUID  `gorm:"type:uuid" json:"live_session_id"`
	StudentID        uuid.UUID  `gorm:"type:uuid" json:"student_id"`
	AttendanceStatus string     `gorm:"size:50" json:"attendance_status"`
	JoinedAt         *time.Time `json:"joined_at,omitempty"`
	LeftAt           *time.Time `json:"left_at,omitempty"`
}
