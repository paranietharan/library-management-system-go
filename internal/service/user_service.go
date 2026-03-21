package service

import (
	"errors"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/repository"
	"library-management-system-go/pkg/utils"
	"strings"
	"time"
)

type UserService interface {
	ListUsers(page, limit int, search string) ([]dto.UserDTO, int64, error)
	GetUser(id uint) (*dto.UserDTO, error)
	CreateUser(req dto.CreateUserRequest) (*dto.UserDTO, error)
	UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserDTO, error)
	DeleteUser(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) ListUsers(page, limit int, search string) ([]dto.UserDTO, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	users, total, err := s.repo.FindAll(page, limit, search)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]dto.UserDTO, 0, len(users))
	for i := range users {
		resp = append(resp, *s.userToDTO(&users[i]))
	}
	return resp, total, nil
}

func (s *userService) GetUser(id uint) (*dto.UserDTO, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.userToDTO(user), nil
}

func (s *userService) CreateUser(req dto.CreateUserRequest) (*dto.UserDTO, error) {
	role := domain.UserRole(strings.ToUpper(strings.TrimSpace(req.Role)))
	if role != domain.RoleAdmin && role != domain.RoleLibrarian && role != domain.RoleTeacher && role != domain.RoleStudent {
		return nil, errors.New("invalid role")
	}

	if role == domain.RoleStudent {
		if req.StudentID == nil || strings.TrimSpace(*req.StudentID) == "" {
			return nil, errors.New("student_id is required for STUDENT role")
		}
		if req.EmployeeID != nil && strings.TrimSpace(*req.EmployeeID) != "" {
			return nil, errors.New("employee_id must be empty for STUDENT role")
		}
	} else {
		if req.EmployeeID == nil || strings.TrimSpace(*req.EmployeeID) == "" {
			return nil, errors.New("employee_id is required for this role")
		}
		if req.StudentID != nil && strings.TrimSpace(*req.StudentID) != "" {
			return nil, errors.New("student_id must be empty for non-STUDENT roles")
		}
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	var dob *time.Time
	if req.DateOfBirth != nil && strings.TrimSpace(*req.DateOfBirth) != "" {
		parsed, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			return nil, errors.New("invalid date_of_birth; expected YYYY-MM-DD")
		}
		dob = &parsed
	}

	user := &domain.User{
		Username:        req.Username,
		Email:           req.Email,
		PasswordHash:    hashedPassword,
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		Role:            role,
		Status:          domain.StatusActive,
		Phone:           req.Phone,
		DateOfBirth:     dob,
		StudentID:       req.StudentID,
		EmployeeID:      req.EmployeeID,
		MaxBooksAllowed: getMaxBooksForRole(role),
	}

	// Ensure uniqueness errors are surfaced cleanly.
	if _, err := s.repo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}
	if _, err := s.repo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return s.userToDTO(user), nil
}

func (s *userService) UpdateUser(id uint, req dto.UpdateUserRequest) (*dto.UserDTO, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}

	if req.DateOfBirth != nil {
		if strings.TrimSpace(*req.DateOfBirth) == "" {
			user.DateOfBirth = nil
		} else {
			parsed, err := time.Parse("2006-01-02", *req.DateOfBirth)
			if err != nil {
				return nil, errors.New("invalid date_of_birth; expected YYYY-MM-DD")
			}
			user.DateOfBirth = &parsed
		}
	}

	if req.Password != nil && strings.TrimSpace(*req.Password) != "" {
		hashedPassword, err := utils.HashPassword(*req.Password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.PasswordHash = hashedPassword
	}

	// Role update includes ID reshaping.
	if req.Role != nil && strings.TrimSpace(*req.Role) != "" {
		newRole := domain.UserRole(strings.ToUpper(strings.TrimSpace(*req.Role)))
		if newRole != domain.RoleAdmin && newRole != domain.RoleLibrarian && newRole != domain.RoleTeacher && newRole != domain.RoleStudent {
			return nil, errors.New("invalid role")
		}
		user.Role = newRole
		user.MaxBooksAllowed = getMaxBooksForRole(newRole)

		// If role changed, require correct IDs.
		if newRole == domain.RoleStudent {
			if req.StudentID == nil || strings.TrimSpace(*req.StudentID) == "" {
				return nil, errors.New("student_id is required for STUDENT role")
			}
			user.StudentID = req.StudentID
			user.EmployeeID = nil
		} else {
			if req.EmployeeID == nil || strings.TrimSpace(*req.EmployeeID) == "" {
				return nil, errors.New("employee_id is required for this role")
			}
			user.EmployeeID = req.EmployeeID
			user.StudentID = nil
		}
	} else {
		// Partial ID updates (only when provided).
		if req.StudentID != nil {
			user.StudentID = req.StudentID
		}
		if req.EmployeeID != nil {
			user.EmployeeID = req.EmployeeID
		}
	}

	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		user.Status = domain.UserStatus(strings.ToUpper(strings.TrimSpace(*req.Status)))
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return s.userToDTO(user), nil
}

func (s *userService) DeleteUser(id uint) error {
	return s.repo.Delete(id)
}

func (s *userService) userToDTO(user *domain.User) *dto.UserDTO {
	return &dto.UserDTO{
		ID:              user.ID,
		Username:        user.Username,
		Email:           user.Email,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		Role:            string(user.Role),
		Status:          string(user.Status),
		Phone:           user.Phone,
		StudentID:       user.StudentID,
		EmployeeID:      user.EmployeeID,
		MaxBooksAllowed: user.MaxBooksAllowed,
	}
}

func getMaxBooksForRole(role domain.UserRole) int {
	switch role {
	case domain.RoleAdmin:
		return 20
	case domain.RoleLibrarian:
		return 15
	case domain.RoleTeacher:
		return 10
	case domain.RoleStudent:
		return 5
	default:
		return 5
	}
}

