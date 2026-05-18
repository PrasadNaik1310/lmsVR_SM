package models

import (
	"time"

	"github.com/google/uuid"
)

type CourseInvite struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	CourseID   uuid.UUID  `gorm:"type:uuid" json:"course_id"`
	InviteCode string     `gorm:"size:100;unique" json:"invite_code"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	CreatedBy  uuid.UUID  `gorm:"type:uuid" json:"created_by"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
}
