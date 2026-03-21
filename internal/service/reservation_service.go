package service

import (
	"errors"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/repository"
	"strings"
	"time"
)

type ReservationService interface {
	CreateReservation(userID uint, role domain.UserRole, req dto.CreateReservationRequest) (*domain.Reservation, error)
	GetReservation(id uint) (*domain.Reservation, error)
	ListReservations(page, limit int, search string) ([]domain.Reservation, int64, error)
	UpdateReservation(userID uint, role domain.UserRole, id uint, req dto.UpdateReservationRequest) (*domain.Reservation, error)
	DeleteReservation(userID uint, role domain.UserRole, id uint) error
}

type reservationService struct {
	repo repository.ReservationRepository
}

func NewReservationService(repo repository.ReservationRepository) ReservationService {
	return &reservationService{repo: repo}
}

func (s *reservationService) CreateReservation(userID uint, _ domain.UserRole, req dto.CreateReservationRequest) (*domain.Reservation, error) {
	if req.ExpiryDate.Before(time.Now()) {
		return nil, errors.New("expiry_date must be in the future")
	}

	reservation := &domain.Reservation{
		UserID:           userID,
		BookID:           req.BookID,
		Status:           domain.ReservationStatusPending,
		ReservationDate: time.Now(),
		ExpiryDate:      req.ExpiryDate,
		FulfilledDate:   nil,
		CancelledDate:   nil,
		LendingID:       nil,
		Notes:           req.Notes,
	}

	if err := s.repo.Create(reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *reservationService) GetReservation(id uint) (*domain.Reservation, error) {
	return s.repo.FindByID(id)
}

func (s *reservationService) ListReservations(page, limit int, search string) ([]domain.Reservation, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.FindAll(page, limit, search)
}

func (s *reservationService) UpdateReservation(userID uint, role domain.UserRole, id uint, req dto.UpdateReservationRequest) (*domain.Reservation, error) {
	reservation, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if reservation.UserID != userID && role != domain.RoleAdmin && role != domain.RoleLibrarian {
		return nil, errors.New("unauthorized to update this reservation")
	}

	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		reservation.Status = domain.ReservationStatus(strings.ToUpper(*req.Status))
	}
	if req.FulfilledDate != nil {
		reservation.FulfilledDate = req.FulfilledDate
	}
	if req.CancelledDate != nil {
		reservation.CancelledDate = req.CancelledDate
	}
	if req.LendingID != nil {
		reservation.LendingID = req.LendingID
	}
	if req.Notes != nil {
		reservation.Notes = req.Notes
	}

	if err := s.repo.Update(reservation); err != nil {
		return nil, err
	}
	return reservation, nil
}

func (s *reservationService) DeleteReservation(userID uint, role domain.UserRole, id uint) error {
	reservation, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if reservation.UserID != userID && role != domain.RoleAdmin && role != domain.RoleLibrarian {
		return errors.New("unauthorized to delete this reservation")
	}

	return s.repo.Delete(id)
}

