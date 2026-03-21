package domain

import "time"

type FineStatus string

const (
	FineStatusPending FineStatus = "PENDING"
	FineStatusPaid    FineStatus = "PAID"
	FineStatusWaived  FineStatus = "WAIVED"
)

type Fine struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	UserID              uint      `gorm:"not null" json:"user_id"`
	LendingID           *uint     `gorm:"default:null" json:"lending_id,omitempty"`
	Amount              float64   `gorm:"not null" json:"amount"`
	Reason              string    `gorm:"not null" json:"reason"`
	Status              FineStatus `gorm:"type:fine_status;not null;default:'PENDING'" json:"status"`
	IssuedDate          time.Time `gorm:"not null" json:"issued_date"`
	DueDate             time.Time `gorm:"not null" json:"due_date"`
	PaidDate            *time.Time `json:"paid_date,omitempty"`
	PaymentMethod      *string    `json:"payment_method,omitempty"`
	PaymentReference   *string    `json:"payment_reference,omitempty"`
	WaivedBy            *uint      `json:"waived_by,omitempty"`
	WaivedDate          *time.Time `json:"waived_date,omitempty"`
	WaiveReason         *string    `json:"waive_reason,omitempty"`
	Notes               *string    `json:"notes,omitempty"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
}

