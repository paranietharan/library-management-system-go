package dto

import "time"

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type CommentResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	BookID    uint      `json:"book_id"`
	Content   string    `json:"content"`
	Username  string    `json:"username,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
