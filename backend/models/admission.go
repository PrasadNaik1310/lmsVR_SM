package models

import (
	"time"

	"github.com/google/uuid"
)

type Enquiry struct {
	ID                 uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	FullName           string    `gorm:"size:255;not null" json:"full_name"`
	Email              string    `gorm:"size:255" json:"email"`
	Phone              string    `gorm:"size:20" json:"phone"`
	InterestedCourseID uuid.UUID `gorm:"type:uuid" json:"interested_course_id"`
	Status             string    `gorm:"size:50" json:"status"`
	Notes              string    `gorm:"type:text" json:"notes"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type Application struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	EnquiryID         uuid.UUID `gorm:"type:uuid" json:"enquiry_id"`
	AppliedCourseID   uuid.UUID `gorm:"type:uuid" json:"applied_course_id"`
	ApplicationStatus string    `gorm:"size:50" json:"application_status"`
	Remarks           string    `gorm:"type:text" json:"remarks"`
	SubmittedAt       time.Time `gorm:"autoCreateTime" json:"submitted_at"`
}
