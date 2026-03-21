package domain

import "time"

type ComplaintStatus string

const (
	ComplaintStatusOpen         ComplaintStatus = "OPEN"
	ComplaintStatusInProgress  ComplaintStatus = "IN_PROGRESS"
	ComplaintStatusResolved    ComplaintStatus = "RESOLVED"
	ComplaintStatusClosed      ComplaintStatus = "CLOSED"
)

type ComplaintPriority string

const (
	ComplaintPriorityLow    ComplaintPriority = "LOW"
	ComplaintPriorityMedium ComplaintPriority = "MEDIUM"
	ComplaintPriorityHigh   ComplaintPriority = "HIGH"
	ComplaintPriorityUrgent ComplaintPriority = "URGENT"
)

type Complaint struct {
	ID               uint               `gorm:"primaryKey" json:"id"`
	UserID          uint               `gorm:"not null" json:"user_id"`
	Subject         string             `gorm:"not null" json:"subject"`
	Description     string             `gorm:"not null" json:"description"`
	Category        *string            `json:"category,omitempty"`
	Priority        ComplaintPriority `gorm:"default:'MEDIUM'" json:"priority"`
	Status          ComplaintStatus   `gorm:"type:complaint_status;default:'OPEN'" json:"status"`
	SubmittedDate   time.Time         `gorm:"not null" json:"submitted_date"`
	AssignedTo      *uint             `json:"assigned_to,omitempty"`
	ResolvedDate    *time.Time        `json:"resolved_date,omitempty"`
	ResolutionNotes *string           `json:"resolution_notes,omitempty"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
}

