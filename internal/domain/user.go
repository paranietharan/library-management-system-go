package domain

import (
	"time"
)

type UserRole string
type UserStatus string

const (
	RoleAdmin     UserRole = "ADMIN"
	RoleLibrarian UserRole = "LIBRARIAN"
	RoleStudent   UserRole = "STUDENT"
	RoleTeacher   UserRole = "TEACHER"
)

const (
	StatusActive    UserStatus = "ACTIVE"
	StatusInactive  UserStatus = "INACTIVE"
	StatusSuspended UserStatus = "SUSPENDED"
)

type User struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Username        string     `gorm:"unique;not null" json:"username"`
	Email           string     `gorm:"unique;not null" json:"email"`
	PasswordHash    string     `gorm:"not null" json:"-"`
	FirstName       string     `gorm:"not null" json:"first_name"`
	LastName        string     `gorm:"not null" json:"last_name"`
	Role            UserRole   `gorm:"type:user_role;not null;default:'STUDENT'" json:"role"`
	Status          UserStatus `gorm:"type:user_status;not null;default:'ACTIVE'" json:"status"`
	Phone           *string    `json:"phone,omitempty"`
	Address         *string    `json:"address,omitempty"`
	DateOfBirth     *time.Time `json:"date_of_birth,omitempty"`
	StudentID       *string    `gorm:"unique" json:"student_id,omitempty"`
	EmployeeID      *string    `gorm:"unique" json:"employee_id,omitempty"`
	MaxBooksAllowed int        `gorm:"default:5" json:"max_books_allowed"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
	LastLogin       *time.Time `json:"last_login,omitempty"`
}
