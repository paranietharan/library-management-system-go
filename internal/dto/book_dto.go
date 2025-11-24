package dto

import "time"

type CreateBookRequest struct {
	Title           string `json:"title" binding:"required"`
	Author          string `json:"author" binding:"required"`
	ISBN            string `json:"isbn" binding:"required"`
	PublicationYear int    `json:"publication_year" binding:"required"`
	Category        string `json:"category" binding:"required"`
	TotalCopies     int    `json:"total_copies" binding:"required,min=1"`
}

type UpdateBookRequest struct {
	Title           string `json:"title"`
	Author          string `json:"author"`
	ISBN            string `json:"isbn"`
	PublicationYear int    `json:"publication_year"`
	Category        string `json:"category"`
	TotalCopies     int    `json:"total_copies"`
	Status          string `json:"status"`
}

type BookResponse struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	ISBN            string    `json:"isbn"`
	PublicationYear int       `json:"publication_year"`
	Category        string    `json:"category"`
	TotalCopies     int       `json:"total_copies"`
	AvailableCopies int       `json:"available_copies"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
