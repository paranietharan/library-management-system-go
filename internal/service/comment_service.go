package service

import (
	"errors"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/repository"
)

type CommentService interface {
	CreateComment(userID, bookID uint, req dto.CreateCommentRequest) (*domain.Comment, error)
	ListComments(bookID uint) ([]domain.Comment, error)
	UpdateComment(userID, commentID uint, userRole domain.UserRole, req dto.UpdateCommentRequest) (*domain.Comment, error)
	DeleteComment(userID, commentID uint, userRole domain.UserRole) error
}

type commentService struct {
	repo repository.CommentRepository
}

func NewCommentService(repo repository.CommentRepository) CommentService {
	return &commentService{repo: repo}
}

func (s *commentService) CreateComment(userID, bookID uint, req dto.CreateCommentRequest) (*domain.Comment, error) {
	comment := &domain.Comment{
		UserID:  userID,
		BookID:  bookID,
		Content: req.Content,
	}

	if err := s.repo.Create(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) ListComments(bookID uint) ([]domain.Comment, error) {
	return s.repo.FindByBookID(bookID)
}

func (s *commentService) UpdateComment(userID, commentID uint, userRole domain.UserRole, req dto.UpdateCommentRequest) (*domain.Comment, error) {
	comment, err := s.repo.FindByID(commentID)
	if err != nil {
		return nil, err
	}

	// Check permissions: Owner, Admin, or Librarian
	if comment.UserID != userID && userRole != domain.RoleAdmin && userRole != domain.RoleLibrarian {
		return nil, errors.New("unauthorized to update this comment")
	}

	if req.Content != "" {
		comment.Content = req.Content
	}

	if err := s.repo.Update(comment); err != nil {
		return nil, err
	}

	return comment, nil
}

func (s *commentService) DeleteComment(userID, commentID uint, userRole domain.UserRole) error {
	comment, err := s.repo.FindByID(commentID)
	if err != nil {
		return err
	}

	// Check permissions: Owner, Admin, or Librarian
	if comment.UserID != userID && userRole != domain.RoleAdmin && userRole != domain.RoleLibrarian {
		return errors.New("unauthorized to delete this comment")
	}

	return s.repo.Delete(commentID)
}
