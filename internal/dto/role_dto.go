package dto

type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description *string `json:"description,omitempty"`
}

type CreateRoleRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
}

