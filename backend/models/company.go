package models

import (
	"time"

	"github.com/google/uuid"
)

type AcademicSession struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name      string    `gorm:"size:100;unique" json:"name"`
	StartDate time.Time `gorm:"not null" json:"start_date"`
	EndDate   time.Time `gorm:"not null" json:"end_date"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
}

type Team struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"size:100;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type TeamMember struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TeamID   uuid.UUID `gorm:"type:uuid" json:"team_id"`
	UserID   uuid.UUID `gorm:"type:uuid" json:"user_id"`
	JoinedAt time.Time `gorm:"autoCreateTime" json:"joined_at"`
}
