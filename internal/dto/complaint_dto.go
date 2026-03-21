package dto

import "time"

type CreateComplaintRequest struct {
	Subject     string  `json:"subject" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Category    *string `json:"category,omitempty"`
	Priority    *string `json:"priority,omitempty"`
}

type UpdateComplaintRequest struct {
	Subject         *string   `json:"subject,omitempty"`
	Description     *string   `json:"description,omitempty"`
	Category        *string   `json:"category,omitempty"`
	Priority        *string   `json:"priority,omitempty"`
	Status          *string   `json:"status,omitempty"`
	AssignedTo      *uint     `json:"assigned_to,omitempty"`
	ResolvedDate    *time.Time `json:"resolved_date,omitempty"`
	ResolutionNotes *string    `json:"resolution_notes,omitempty"`
}

type ComplaintResponse struct {
	ID               uint      `json:"id"`
	UserID          uint      `json:"user_id"`
	Subject         string    `json:"subject"`
	Description     string    `json:"description"`
	Category        *string   `json:"category,omitempty"`
	Priority        string    `json:"priority"`
	Status          string    `json:"status"`
	SubmittedDate   time.Time `json:"submitted_date"`
	AssignedTo      *uint     `json:"assigned_to,omitempty"`
	ResolvedDate    *time.Time `json:"resolved_date,omitempty"`
	ResolutionNotes *string   `json:"resolution_notes,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

