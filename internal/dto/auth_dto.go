package dto

type RegisterRequest struct {
	Username    string  `json:"username" binding:"required,min=3,max=50"`
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	FirstName   string  `json:"first_name" binding:"required"`
	LastName    string  `json:"last_name" binding:"required"`
	Phone       *string `json:"phone,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	StudentID   *string `json:"student_id,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string   `json:"token"`
	User      *UserDTO `json:"user"`
	ExpiresAt int64    `json:"expires_at"`
}

type UserDTO struct {
	ID              uint    `json:"id"`
	Username        string  `json:"username"`
	Email           string  `json:"email"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Role            string  `json:"role"`
	Status          string  `json:"status"`
	Phone           *string `json:"phone,omitempty"`
	StudentID       *string `json:"student_id,omitempty"`
	EmployeeID      *string `json:"employee_id,omitempty"`
	MaxBooksAllowed int     `json:"max_books_allowed"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}
