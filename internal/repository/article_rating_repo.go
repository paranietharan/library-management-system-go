package repository

import (
	"library-management-system-go/internal/domain"

	"gorm.io/gorm"
)

type ArticleRatingRepository interface {
	Create(rating *domain.ArticleRating) error
	FindByID(id uint) (*domain.ArticleRating, error)
	FindByArticleID(articleID uint) ([]domain.ArticleRating, error)
	Update(rating *domain.ArticleRating) error
	Delete(id uint) error
}

type articleRatingRepository struct {
	db *gorm.DB
}

func NewArticleRatingRepository(db *gorm.DB) ArticleRatingRepository {
	return &articleRatingRepository{db: db}
}

func (r *articleRatingRepository) Create(rating *domain.ArticleRating) error {
	return r.db.Create(rating).Error
}

func (r *articleRatingRepository) FindByID(id uint) (*domain.ArticleRating, error) {
	var rating domain.ArticleRating
	if err := r.db.First(&rating, id).Error; err != nil {
		return nil, err
	}
	return &rating, nil
}

func (r *articleRatingRepository) FindByArticleID(articleID uint) ([]domain.ArticleRating, error) {
	var ratings []domain.ArticleRating
	if err := r.db.Where("article_id = ?", articleID).Find(&ratings).Error; err != nil {
		return nil, err
	}
	return ratings, nil
}

func (r *articleRatingRepository) Update(rating *domain.ArticleRating) error {
	return r.db.Save(rating).Error
}

func (r *articleRatingRepository) Delete(id uint) error {
	return r.db.Delete(&domain.ArticleRating{}, id).Error
}

