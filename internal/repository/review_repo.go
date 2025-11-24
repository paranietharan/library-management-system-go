package repository

import (
	"library-management-system-go/internal/domain"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review *domain.Review) error
	FindByID(id uint) (*domain.Review, error)
	FindByBookID(bookID uint) ([]domain.Review, error)
	Update(review *domain.Review) error
	Delete(id uint) error
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *domain.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) FindByID(id uint) (*domain.Review, error) {
	var review domain.Review
	err := r.db.Preload("User").First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) FindByBookID(bookID uint) ([]domain.Review, error) {
	var reviews []domain.Review
	err := r.db.Where("book_id = ?", bookID).Preload("User").Find(&reviews).Error
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *reviewRepository) Update(review *domain.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Review{}, id).Error
}
