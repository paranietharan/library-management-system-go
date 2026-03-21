package dto

type CreateUserRequest struct {
	Username   string  `json:"username" binding:"required,min=3,max=50"`
	Email      string  `json:"email" binding:"required,email"`
	Password   string  `json:"password" binding:"required,min=6"`
	FirstName  string  `json:"first_name" binding:"required"`
	LastName   string  `json:"last_name" binding:"required"`
	Role       string  `json:"role" binding:"required"`
	Phone      *string `json:"phone,omitempty"`
	StudentID  *string `json:"student_id,omitempty"`
	EmployeeID *string `json:"employee_id,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
}

type UpdateUserRequest struct {
	FirstName  *string `json:"first_name,omitempty"`
	LastName   *string `json:"last_name,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	StudentID  *string `json:"student_id,omitempty"`
	EmployeeID *string `json:"employee_id,omitempty"`
	Role       *string `json:"role,omitempty"`
	Status     *string `json:"status,omitempty"`
	Password   *string `json:"password,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
}

