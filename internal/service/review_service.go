package service

import (
	"errors"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/repository"
)

type ReviewService interface {
	CreateReview(userID, bookID uint, req dto.CreateReviewRequest) (*domain.Review, error)
	ListReviews(bookID uint) ([]domain.Review, error)
	UpdateReview(userID, reviewID uint, userRole domain.UserRole, req dto.UpdateReviewRequest) (*domain.Review, error)
	DeleteReview(userID, reviewID uint, userRole domain.UserRole) error
}

type reviewService struct {
	repo repository.ReviewRepository
}

func NewReviewService(repo repository.ReviewRepository) ReviewService {
	return &reviewService{repo: repo}
}

func (s *reviewService) CreateReview(userID, bookID uint, req dto.CreateReviewRequest) (*domain.Review, error) {
	review := &domain.Review{
		UserID:  userID,
		BookID:  bookID,
		Rating:  req.Rating,
		Comment: req.Comment,
	}

	if err := s.repo.Create(review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewService) ListReviews(bookID uint) ([]domain.Review, error) {
	return s.repo.FindByBookID(bookID)
}

func (s *reviewService) UpdateReview(userID, reviewID uint, userRole domain.UserRole, req dto.UpdateReviewRequest) (*domain.Review, error) {
	review, err := s.repo.FindByID(reviewID)
	if err != nil {
		return nil, err
	}

	// Check permissions: Owner, Admin, or Librarian
	if review.UserID != userID && userRole != domain.RoleAdmin && userRole != domain.RoleLibrarian {
		return nil, errors.New("unauthorized to update this review")
	}

	if req.Rating != 0 {
		review.Rating = req.Rating
	}
	if req.Comment != "" {
		review.Comment = req.Comment
	}

	if err := s.repo.Update(review); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *reviewService) DeleteReview(userID, reviewID uint, userRole domain.UserRole) error {
	review, err := s.repo.FindByID(reviewID)
	if err != nil {
		return err
	}

	// Check permissions: Owner, Admin, or Librarian
	if review.UserID != userID && userRole != domain.RoleAdmin && userRole != domain.RoleLibrarian {
		return errors.New("unauthorized to delete this review")
	}

	return s.repo.Delete(reviewID)
}
