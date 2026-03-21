package domain

import "time"

type LendingStatus string

const (
	LendingStatusActive  LendingStatus = "ACTIVE"
	LendingStatusReturned LendingStatus = "RETURNED"
	LendingStatusOverdue LendingStatus = "OVERDUE"
	LendingStatusLost    LendingStatus = "LOST"
)

type Lending struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	UserID        uint          `gorm:"not null" json:"user_id"`
	BookID        uint          `gorm:"not null" json:"book_id"`
	Status        LendingStatus `gorm:"type:lending_status;not null;default:'ACTIVE'" json:"status"`
	IssueDate     time.Time     `gorm:"not null;default:CURRENT_TIMESTAMP" json:"issue_date"`
	DueDate       time.Time     `gorm:"not null" json:"due_date"`
	ReturnDate    *time.Time    `json:"return_date,omitempty"`
	RenewalCount  int           `gorm:"default:0" json:"renewal_count"`
	MaxRenewals   int           `gorm:"default:2" json:"max_renewals"`
	Notes         *string       `json:"notes,omitempty"`
	IssuedBy      *uint         `json:"issued_by,omitempty"`
	ReturnedTo    *uint         `json:"returned_to,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

