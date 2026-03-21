package dto

import "time"

type CreateLendingRequest struct {
	BookID      uint      `json:"book_id" binding:"required"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	Notes       *string   `json:"notes,omitempty"`
	MaxRenewals *int      `json:"max_renewals,omitempty"`
}

type UpdateLendingRequest struct {
	Status       *string   `json:"status,omitempty"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	ReturnDate   *time.Time `json:"return_date,omitempty"`
	RenewalCount *int      `json:"renewal_count,omitempty"`
	MaxRenewals  *int      `json:"max_renewals,omitempty"`
	Notes        *string   `json:"notes,omitempty"`
}

type LendingResponse struct {
	ID           uint       `json:"id"`
	UserID       uint       `json:"user_id"`
	BookID       uint       `json:"book_id"`
	Status       string     `json:"status"`
	IssueDate    time.Time  `json:"issue_date"`
	DueDate      time.Time  `json:"due_date"`
	ReturnDate   *time.Time `json:"return_date,omitempty"`
	RenewalCount int        `json:"renewal_count"`
	MaxRenewals  int        `json:"max_renewals"`
	Notes        *string    `json:"notes,omitempty"`
	IssuedBy     *uint      `json:"issued_by,omitempty"`
	ReturnedTo   *uint      `json:"returned_to,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

