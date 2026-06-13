package models

import (
	"time"

	"github.com/google/uuid"
)

type Student struct {
	ID     uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID uuid.UUID `gorm:"type:uuid;unique;not null" json:"user_id"`
	//EnrollmentNumber string     `gorm:"size:200;unique" json:"enrollment_number"`
	EnrollmentNumber string     `gorm:"unique" json:"enrollment_number"`
	DateOfBirth      *time.Time `json:"date_of_birth,omitempty"`
	Address          string     `gorm:"type:text" json:"address"`
	AdmissionDate    *time.Time `json:"admission_date,omitempty" gorm:"not null"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

type Teacher struct {
	ID             uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         uuid.UUID  `gorm:"type:uuid;unique" json:"user_id"`
	Specialization string     `gorm:"size:255" json:"specialization"`
	Bio            string     `gorm:"type:text" json:"bio"`
	JoiningDate    *time.Time `json:"joining_date,omitempty"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
