package service

import (
	"errors"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/repository"
	"time"
)

type FineService interface {
	CreateFine(userID uint, role domain.UserRole, req dto.CreateFineRequest) (*domain.Fine, error)
	GetFine(id uint) (*domain.Fine, error)
	ListFines(page, limit int, search string) ([]domain.Fine, int64, error)
}

type fineService struct {
	fineRepo    repository.FineRepository
	lendingRepo repository.LendingRepository
}

func NewFineService(fineRepo repository.FineRepository, lendingRepo repository.LendingRepository) FineService {
	return &fineService{
		fineRepo:    fineRepo,
		lendingRepo: lendingRepo,
	}
}

func (s *fineService) CreateFine(userID uint, role domain.UserRole, req dto.CreateFineRequest) (*domain.Fine, error) {
	lending, err := s.lendingRepo.FindByID(req.LendingID)
	if err != nil {
		return nil, errors.New("lending not found")
	}

	if lending.UserID != userID && role != domain.RoleAdmin && role != domain.RoleLibrarian {
		return nil, errors.New("insufficient permissions to create this fine")
	}

	fine := &domain.Fine{
		UserID:       lending.UserID,
		LendingID:    &lending.ID,
		Amount:       req.Amount,
		Reason:       req.Reason,
		Status:       domain.FineStatusPending,
		IssuedDate:   time.Now(),
		DueDate:      req.DueDate,
		PaidDate:     nil,
		Notes:        req.Notes,
	}

	if err := s.fineRepo.Create(fine); err != nil {
		return nil, err
	}

	return fine, nil
}

func (s *fineService) GetFine(id uint) (*domain.Fine, error) {
	return s.fineRepo.FindByID(id)
}

func (s *fineService) ListFines(page, limit int, search string) ([]domain.Fine, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.fineRepo.FindAll(page, limit, search)
}

