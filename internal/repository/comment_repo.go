package repository

import (
	"library-management-system-go/internal/domain"

	"gorm.io/gorm"
)

type CommentRepository interface {
	Create(comment *domain.Comment) error
	FindByID(id uint) (*domain.Comment, error)
	FindByBookID(bookID uint) ([]domain.Comment, error)
	Update(comment *domain.Comment) error
	Delete(id uint) error
}

type commentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &commentRepository{db: db}
}

func (r *commentRepository) Create(comment *domain.Comment) error {
	return r.db.Create(comment).Error
}

func (r *commentRepository) FindByID(id uint) (*domain.Comment, error) {
	var comment domain.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *commentRepository) FindByBookID(bookID uint) ([]domain.Comment, error) {
	var comments []domain.Comment
	err := r.db.Where("book_id = ?", bookID).Preload("User").Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *commentRepository) Update(comment *domain.Comment) error {
	return r.db.Save(comment).Error
}

func (r *commentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Comment{}, id).Error
}
