package service

import (
	"errors"
	"library-management-system-go/internal/dto"
	"strings"
)

type RoleService interface {
	ListRoles(page, limit int, search string) ([]dto.RoleResponse, int64, error)
	GetRole(id uint) (*dto.RoleResponse, error)
	CreateRole(req dto.CreateRoleRequest) (*dto.RoleResponse, error)
}

type roleService struct {
	roles map[uint]dto.RoleResponse
}

func NewRoleService() RoleService {
	// Roles are fixed to the `user_role` PostgreSQL enum.
	// IDs are stable for API responses.
	description := func(s string) *string { return &s }
	return &roleService{
		roles: map[uint]dto.RoleResponse{
			1: {ID: 1, Name: "ADMIN", Description: description("Has all permissions")},
			2: {ID: 2, Name: "LIBRARIAN", Description: description("Manages library inventory and users")},
			3: {ID: 3, Name: "TEACHER", Description: description("Reviews articles and manages learning workflows")},
			4: {ID: 4, Name: "STUDENT", Description: description("Submits articles and uses lending features")},
		},
	}
}

func (s *roleService) ListRoles(page, limit int, search string) ([]dto.RoleResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	search = strings.TrimSpace(search)
	all := make([]dto.RoleResponse, 0, len(s.roles))
	for _, r := range s.roles {
		if search == "" || strings.Contains(strings.ToUpper(r.Name), strings.ToUpper(search)) {
			all = append(all, r)
		}
	}

	total := int64(len(all))
	start := (page - 1) * limit
	if start >= len(all) {
		return []dto.RoleResponse{}, total, nil
	}

	end := start + limit
	if end > len(all) {
		end = len(all)
	}

	return all[start:end], total, nil
}

func (s *roleService) GetRole(id uint) (*dto.RoleResponse, error) {
	r, ok := s.roles[id]
	if !ok {
		return nil, errors.New("role not found")
	}
	return &r, nil
}

func (s *roleService) CreateRole(req dto.CreateRoleRequest) (*dto.RoleResponse, error) {
	name := strings.ToUpper(strings.TrimSpace(req.Name))
	for _, r := range s.roles {
		if r.Name == name {
			// Idempotent: role already exists in the enum.
			return &r, nil
		}
	}
	return nil, errors.New("roles are fixed by system enum; cannot create unknown role")
}

