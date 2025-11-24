package dto

import "time"

type CreateReviewRequest struct {
	Rating  int    `json:"rating" binding:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating" binding:"min=1,max=5"`
	Comment string `json:"comment"`
}

type ReviewResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	BookID    uint      `json:"book_id"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
