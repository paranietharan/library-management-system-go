package dto

import "time"

type CreateReservationRequest struct {
	BookID     uint      `json:"book_id" binding:"required"`
	ExpiryDate time.Time `json:"expiry_date" binding:"required"`
	Notes      *string   `json:"notes,omitempty"`
}

type UpdateReservationRequest struct {
	Status        *string   `json:"status,omitempty"`
	FulfilledDate *time.Time `json:"fulfilled_date,omitempty"`
	CancelledDate *time.Time `json:"cancelled_date,omitempty"`
	LendingID     *uint     `json:"lending_id,omitempty"`
	Notes         *string   `json:"notes,omitempty"`
}

type ReservationResponse struct {
	ID               uint       `json:"id"`
	UserID          uint       `json:"user_id"`
	BookID          uint       `json:"book_id"`
	Status          string     `json:"status"`
	ReservationDate time.Time  `json:"reservation_date"`
	ExpiryDate      time.Time  `json:"expiry_date"`
	FulfilledDate   *time.Time `json:"fulfilled_date,omitempty"`
	CancelledDate   *time.Time `json:"cancelled_date,omitempty"`
	LendingID       *uint      `json:"lending_id,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

