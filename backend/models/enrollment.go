package models

import (
	"time"

	"github.com/google/uuid"
)

type Enrollment struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID        uuid.UUID `gorm:"type:uuid" json:"student_id"`
	CourseID         uuid.UUID `gorm:"type:uuid" json:"course_id"`
	EnrollmentStatus string    `gorm:"size:50" json:"enrollment_status"`
	JoinedViaLink    bool      `gorm:"default:false" json:"joined_via_link"`
	JoinedAt         time.Time `gorm:"autoCreateTime" json:"joined_at"`
}

type Membership struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID        uuid.UUID  `gorm:"type:uuid;unique" json:"student_id"`
	CourseID         uuid.UUID  `gorm:"type:uuid" json:"course_id"`
	TotalFee         float64    `gorm:"type:numeric(10,2);not null" json:"total_fee"`
	AmountPaid       float64    `gorm:"type:numeric(10,2);default:0" json:"amount_paid"`
	DueAmount        float64    `gorm:"type:numeric(10,2);default:0" json:"due_amount"`
	GraceExpiryDate  *time.Time `json:"grace_expiry_date,omitempty"`
	MembershipStatus string     `gorm:"size:50" json:"membership_status"`
	CreatedAt        time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

type Payment struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	MembershipID         uuid.UUID  `gorm:"type:uuid" json:"membership_id"`
	StudentID            uuid.UUID  `gorm:"type:uuid" json:"student_id"`
	Amount               float64    `gorm:"type:numeric(10,2);check:amount > 0" json:"amount"`
	PaymentMethod        string     `gorm:"size:50" json:"payment_method"`
	PaymentStatus        string     `gorm:"size:50" json:"payment_status"`
	TransactionReference string     `gorm:"size:255;unique" json:"transaction_reference"`
	PaidAt               *time.Time `json:"paid_at,omitempty"`
}
