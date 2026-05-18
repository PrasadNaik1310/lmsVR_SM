package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	FirstName    string     `gorm:"size:100;not null" json:"first_name"`
	LastName     string     `gorm:"size:100" json:"last_name"`
	Email        string     `gorm:"size:255;unique;not null" json:"email"`
	PasswordHash string     `gorm:"type:text;not null" json:"-"`
	Phone        string     `gorm:"size:20;unique" json:"phone"`
	RoleID       uuid.UUID  `gorm:"type:uuid" json:"role_id"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"size:50;unique;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
}

type Permission struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"size:100;unique;not null" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
}

type RolePermission struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	RoleID       uuid.UUID `gorm:"type:uuid" json:"role_id"`
	PermissionID uuid.UUID `gorm:"type:uuid" json:"permission_id"`
}
