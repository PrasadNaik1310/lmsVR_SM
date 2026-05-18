package models

import (
	"time"

	"github.com/google/uuid"
)

type CoursePlanner struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CourseID         uuid.UUID `gorm:"type:uuid" json:"course_id"`
	LessonID         uuid.UUID `gorm:"type:uuid" json:"lesson_id"`
	PlannedDate      time.Time `gorm:"not null" json:"planned_date"`
	PlannedStartTime string    `gorm:"type:time" json:"planned_start_time"`
	PlannedEndTime   string    `gorm:"type:time" json:"planned_end_time"`
	RecurrenceType   string    `gorm:"size:50" json:"recurrence_type"`
	CreatedBy        uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type CourseLog struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	PlannerID        uuid.UUID  `gorm:"type:uuid" json:"planner_id"`
	CompletionStatus string     `gorm:"size:50" json:"completion_status"`
	Remarks          string     `gorm:"type:text" json:"remarks"`
	CompletedAt      *time.Time `json:"completed_at,omitempty"`
	RecordedBy       uuid.UUID  `gorm:"type:uuid" json:"recorded_by"`
}
