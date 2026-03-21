package service

import (
	"errors"
	"library-management-system-go/internal/config"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/repository"
	"library-management-system-go/pkg/utils"
	"time"

	"gorm.io/gorm"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) (*dto.UserDTO, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetUserByID(id uint) (*dto.UserDTO, error)
	ChangePassword(userID uint, req *dto.ChangePasswordRequest) error
	ValidateToken(token string) (*utils.JWTClaims, error)
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		cfg:      cfg,
	}
}

func (s *authService) Register(req *dto.RegisterRequest) (*dto.UserDTO, error) {
	if req.Role != "STUDENT" && req.Role != "TEACHER" {
		return nil, errors.New("role must be STUDENT or TEACHER")
	}

	if req.Role == "STUDENT" && (req.StudentID == nil || *req.StudentID == "") {
		return nil, errors.New("student_id is required for student")
	}

	if req.Role == "TEACHER" && (req.EmployeeID == nil || *req.EmployeeID == "") {
		return nil, errors.New("employee_id is required for teacher")
	}

	// Check if username already exists
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Parse date of birth if provided
	var dob *time.Time
	if req.DateOfBirth != nil {
		parsedDOB, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err == nil {
			dob = &parsedDOB
		}
	}

	// Create user
	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		DateOfBirth:  dob,
		StudentID:    req.StudentID,
		EmployeeID:   req.EmployeeID,
		Role:         domain.UserRole(req.Role),
		Status:       domain.StatusInactive,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	return s.userToDTO(user), nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Username not found in the database")
		}
		return nil, errors.New("failed to find user")
	}

	if user.Status != domain.StatusActive {
		return nil, errors.New("account is not active")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	token, expiresAt, err := utils.GenerateToken(
		user.ID,
		user.Username,
		user.Email,
		string(user.Role),
		s.cfg.JWT.Secret,
		s.cfg.JWT.ExpiryHours,
	)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

	return &dto.LoginResponse{
		Token:     token,
		User:      s.userToDTO(user),
		ExpiresAt: expiresAt,
	}, nil
}

func (s *authService) GetUserByID(id uint) (*dto.UserDTO, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return s.userToDTO(user), nil
}

func (s *authService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !utils.CheckPassword(req.OldPassword, user.PasswordHash) {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.PasswordHash = hashedPassword
	if err := s.userRepo.Update(user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (s *authService) ValidateToken(token string) (*utils.JWTClaims, error) {
	return utils.ValidateToken(token, s.cfg.JWT.Secret)
}

func (s *authService) userToDTO(user *domain.User) *dto.UserDTO {
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
