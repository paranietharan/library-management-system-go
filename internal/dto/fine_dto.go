package dto

import "time"

type CreateFineRequest struct {
	LendingID uint    `json:"lending_id" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,gt=0"`
	Reason    string  `json:"reason" binding:"required"`
	DueDate   time.Time `json:"due_date" binding:"required"`
	Notes     *string `json:"notes,omitempty"`
}

type FineResponse struct {
	ID                uint       `json:"id"`
	UserID            uint       `json:"user_id"`
	LendingID         *uint      `json:"lending_id,omitempty"`
	Amount            float64    `json:"amount"`
	Reason            string     `json:"reason"`
	Status            string     `json:"status"`
	IssuedDate        time.Time  `json:"issued_date"`
	DueDate           time.Time  `json:"due_date"`
	PaidDate          *time.Time `json:"paid_date,omitempty"`
	PaymentMethod     *string    `json:"payment_method,omitempty"`
	PaymentReference  *string    `json:"payment_reference,omitempty"`
	WaivedBy          *uint      `json:"waived_by,omitempty"`
	WaivedDate        *time.Time `json:"waived_date,omitempty"`
	WaiveReason       *string    `json:"waive_reason,omitempty"`
	Notes             *string    `json:"notes,omitempty"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
}

