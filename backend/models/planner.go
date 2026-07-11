package models

import (
	"time"

	"github.com/google/uuid"
)

type CourseSchedule struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CourseID         uuid.UUID `gorm:"type:uuid" json:"course_id"`
	LessonID         uuid.UUID `gorm:"type:uuid" json:"lesson_id"`
	TeacherID        uuid.UUID `gorm:"type:uuid" json:"teacher_id"`
	PlannedDate      time.Time `gorm:"not null" json:"planned_date"`
	PlannedStartTime string    `gorm:"type:time" json:"planned_start_time"`
	PlannedEndTime   string    `gorm:"type:time" json:"planned_end_time"`
	//RecurrenceType   string    `gorm:"size:50" json:"recurrence_type"`
	//NOTE: uncomment the below if admin coordinator role is requried
	//Assuming teacher is creating scheudles for now .

	Status    string    `gorm:"size:30;default:'SCHEDULED'" json:"status"`
	CreatedBy uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CourseLog struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	ScheduleID       uuid.UUID `gorm:"type:uuid;uniqueIndex" json:"schedule_id"`
	ConductedDate    time.Time `gorm:"not null" json:"conducted_date"`
	CompletionStatus string    `gorm:"size:50;not null" json:"completion_status"`
	Remarks          string    `gorm:"type:text" json:"remarks"`
	Homework         string    `gorm:"type:text" json:"homework"`
	NextTopic        string    `gorm:"type:text" json:"next_topic"`
	RecordedBy       uuid.UUID `gorm:"type:uuid" json:"recorded_by"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
