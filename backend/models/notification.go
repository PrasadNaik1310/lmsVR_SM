package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid" json:"user_id"`
	Type      string    `gorm:"size:100" json:"type"`
	Title     string    `gorm:"size:255" json:"title"`
	Message   string    `gorm:"type:text" json:"message"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	SentEmail bool      `gorm:"default:false" json:"sent_email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
