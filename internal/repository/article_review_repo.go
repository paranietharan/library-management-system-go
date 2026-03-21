package repository

import (
	"library-management-system-go/internal/domain"

	"gorm.io/gorm"
)

type ArticleReviewRepository interface {
	Create(review *domain.ArticleReview) error
	FindByID(id uint) (*domain.ArticleReview, error)
	FindAll(page, limit int, search string) ([]domain.ArticleReview, int64, error)
	Update(review *domain.ArticleReview) error
}

type articleReviewRepository struct {
	db *gorm.DB
}

func NewArticleReviewRepository(db *gorm.DB) ArticleReviewRepository {
	return &articleReviewRepository{db: db}
}

func (r *articleReviewRepository) Create(review *domain.ArticleReview) error {
	return r.db.Create(review).Error
}

func (r *articleReviewRepository) FindByID(id uint) (*domain.ArticleReview, error) {
	var review domain.ArticleReview
	if err := r.db.First(&review, id).Error; err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *articleReviewRepository) FindAll(page, limit int, search string) ([]domain.ArticleReview, int64, error) {
	var reviews []domain.ArticleReview
	var total int64

	query := r.db.Model(&domain.ArticleReview{})
	if search != "" {
		// Simple search: feedback match. (Full search could join `articles`.)
		like := "%" + search + "%"
		query = query.Where("feedback ILIKE ?", like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&reviews).Error; err != nil {
		return nil, 0, err
	}

	return reviews, total, nil
}

func (r *articleReviewRepository) Update(review *domain.ArticleReview) error {
	return r.db.Save(review).Error
}

