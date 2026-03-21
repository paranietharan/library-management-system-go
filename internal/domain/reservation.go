package domain

import "time"

type ReservationStatus string

const (
	ReservationStatusPending   ReservationStatus = "PENDING"
	ReservationStatusFulfilled ReservationStatus = "FULFILLED"
	ReservationStatusCancelled ReservationStatus = "CANCELLED"
	ReservationStatusExpired   ReservationStatus = "EXPIRED"
)

type Reservation struct {
	ID               uint                 `gorm:"primaryKey" json:"id"`
	UserID          uint                 `gorm:"not null" json:"user_id"`
	BookID          uint                 `gorm:"not null" json:"book_id"`
	Status          ReservationStatus   `gorm:"type:reservation_status;not null;default:'PENDING'" json:"status"`
	ReservationDate time.Time           `gorm:"not null;default:CURRENT_TIMESTAMP" json:"reservation_date"`
	ExpiryDate      time.Time           `gorm:"not null" json:"expiry_date"`
	FulfilledDate   *time.Time          `json:"fulfilled_date,omitempty"`
	CancelledDate   *time.Time          `json:"cancelled_date,omitempty"`
	LendingID       *uint               `json:"lending_id,omitempty"`
	Notes           *string            `json:"notes,omitempty"`
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

