package domain

import (
	"time"
)

type BookStatus string

const (
	BookStatusAvailable   BookStatus = "AVAILABLE"
	BookStatusOutOfStock  BookStatus = "OUT_OF_STOCK"
	BookStatusMaintenance BookStatus = "MAINTENANCE"
)

type Book struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Title           string     `gorm:"not null" json:"title"`
	Author          string     `gorm:"not null" json:"author"`
	ISBN            string     `gorm:"unique;not null" json:"isbn"`
	PublicationYear int        `json:"publication_year"`
	Category        string     `json:"category"`
	TotalCopies     int        `gorm:"default:1" json:"total_copies"`
	AvailableCopies int        `gorm:"default:1" json:"available_copies"`
	Status          BookStatus `gorm:"type:book_status;default:'AVAILABLE'" json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
