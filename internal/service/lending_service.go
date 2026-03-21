package service

import (
	"errors"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/repository"
	"strings"
	"time"
)

type LendingService interface {
	CreateLending(userID uint, role domain.UserRole, req dto.CreateLendingRequest) (*domain.Lending, error)
	GetLending(id uint) (*domain.Lending, error)
	ListLendings(page, limit int, search string) ([]domain.Lending, int64, error)
	UpdateLending(userID uint, role domain.UserRole, id uint, req dto.UpdateLendingRequest) (*domain.Lending, error)
	DeleteLending(userID uint, role domain.UserRole, id uint) error
}

type lendingService struct {
	repo repository.LendingRepository
}

func NewLendingService(repo repository.LendingRepository) LendingService {
	return &lendingService{repo: repo}
}

func (s *lendingService) CreateLending(userID uint, _ domain.UserRole, req dto.CreateLendingRequest) (*domain.Lending, error) {
	if req.DueDate.Before(time.Now()) {
		// Not strictly required by DB (it checks due_date > issue_date), but gives better feedback.
		// We'll still allow due_date in future only.
		return nil, errors.New("due_date must be in the future")
	}

	lending := &domain.Lending{
		UserID:       userID,
		BookID:       req.BookID,
		Status:       domain.LendingStatusActive,
		IssueDate:    time.Now(),
		DueDate:      req.DueDate,
		ReturnDate:   nil,
		RenewalCount: 0,
		MaxRenewals:  2,
		Notes:        req.Notes,
		IssuedBy:     &userID,
		ReturnedTo:   nil,
	}

	if req.MaxRenewals != nil {
		lending.MaxRenewals = *req.MaxRenewals
	}

	if err := s.repo.Create(lending); err != nil {
		return nil, err
	}

	return lending, nil
}

func (s *lendingService) GetLending(id uint) (*domain.Lending, error) {
	return s.repo.FindByID(id)
}

func (s *lendingService) ListLendings(page, limit int, search string) ([]domain.Lending, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.FindAll(page, limit, search)
}

func (s *lendingService) UpdateLending(userID uint, role domain.UserRole, id uint, req dto.UpdateLendingRequest) (*domain.Lending, error) {
	lending, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if lending.UserID != userID && role != domain.RoleAdmin && role != domain.RoleLibrarian {
		return nil, errors.New("unauthorized to update this lending")
	}

	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		lending.Status = domain.LendingStatus(strings.ToUpper(*req.Status))
	}
	if req.DueDate != nil {
		lending.DueDate = *req.DueDate
	}
	if req.ReturnDate != nil {
		lending.ReturnDate = req.ReturnDate
	}
	if req.RenewalCount != nil {
		lending.RenewalCount = *req.RenewalCount
	}
	if req.MaxRenewals != nil {
		lending.MaxRenewals = *req.MaxRenewals
	}
	if req.Notes != nil {
		lending.Notes = req.Notes
	}

	if err := s.repo.Update(lending); err != nil {
		return nil, err
	}

	return lending, nil
}

func (s *lendingService) DeleteLending(userID uint, role domain.UserRole, id uint) error {
	lending, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if lending.UserID != userID && role != domain.RoleAdmin && role != domain.RoleLibrarian {
		return errors.New("unauthorized to delete this lending")
	}

	return s.repo.Delete(id)
}

