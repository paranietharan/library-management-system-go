package service

import (
	"errors"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/repository"
	"strings"
	"time"
)

type ComplaintService interface {
	CreateComplaint(userID uint, req dto.CreateComplaintRequest) (*domain.Complaint, error)
	GetComplaint(id uint) (*domain.Complaint, error)
	ListComplaints(page, limit int, search string) ([]domain.Complaint, int64, error)
	UpdateComplaint(userID uint, role domain.UserRole, id uint, req dto.UpdateComplaintRequest) (*domain.Complaint, error)
	DeleteComplaint(userID uint, role domain.UserRole, id uint) error
}

type complaintService struct {
	repo repository.ComplaintRepository
}

func NewComplaintService(repo repository.ComplaintRepository) ComplaintService {
	return &complaintService{repo: repo}
}

func (s *complaintService) CreateComplaint(userID uint, req dto.CreateComplaintRequest) (*domain.Complaint, error) {
	priority := domain.ComplaintPriorityMedium
	if req.Priority != nil && strings.TrimSpace(*req.Priority) != "" {
		priority = domain.ComplaintPriority(strings.ToUpper(*req.Priority))
	}

	complaint := &domain.Complaint{
		UserID:         userID,
		Subject:        req.Subject,
		Description:    req.Description,
		Category:       req.Category,
		Priority:       priority,
		Status:         domain.ComplaintStatusOpen,
		SubmittedDate:  time.Now(),
		AssignedTo:     nil,
		ResolvedDate:   nil,
		ResolutionNotes: nil,
	}

	if err := s.repo.Create(complaint); err != nil {
		return nil, err
	}
	return complaint, nil
}

func (s *complaintService) GetComplaint(id uint) (*domain.Complaint, error) {
	return s.repo.FindByID(id)
}

func (s *complaintService) ListComplaints(page, limit int, search string) ([]domain.Complaint, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.FindAll(page, limit, search)
}

func (s *complaintService) UpdateComplaint(userID uint, role domain.UserRole, id uint, req dto.UpdateComplaintRequest) (*domain.Complaint, error) {
	complaint, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if complaint.UserID != userID && role != domain.RoleAdmin && role != domain.RoleLibrarian {
		return nil, errors.New("unauthorized to update this complaint")
	}

	if req.Subject != nil {
		complaint.Subject = *req.Subject
	}
	if req.Description != nil {
		complaint.Description = *req.Description
	}
	if req.Category != nil {
		complaint.Category = req.Category
	}
	if req.Priority != nil && strings.TrimSpace(*req.Priority) != "" {
		complaint.Priority = domain.ComplaintPriority(strings.ToUpper(*req.Priority))
	}
	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		complaint.Status = domain.ComplaintStatus(strings.ToUpper(*req.Status))
	}
	if req.AssignedTo != nil {
		complaint.AssignedTo = req.AssignedTo
	}
	if req.ResolvedDate != nil {
		complaint.ResolvedDate = req.ResolvedDate
	}
	if req.ResolutionNotes != nil {
		complaint.ResolutionNotes = req.ResolutionNotes
	}

	if err := s.repo.Update(complaint); err != nil {
		return nil, err
	}

	return complaint, nil
}

func (s *complaintService) DeleteComplaint(userID uint, role domain.UserRole, id uint) error {
	complaint, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	if complaint.UserID != userID && role != domain.RoleAdmin && role != domain.RoleLibrarian {
		return errors.New("unauthorized to delete this complaint")
	}

	return s.repo.Delete(id)
}

