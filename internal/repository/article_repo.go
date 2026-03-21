package repository

import (
	"library-management-system-go/internal/domain"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	Create(article *domain.Article) error
	FindByID(id uint) (*domain.Article, error)
	FindAll(page, limit int, search string) ([]domain.Article, int64, error)
	Update(article *domain.Article) error
	Delete(id uint) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Create(article *domain.Article) error {
	return r.db.Create(article).Error
}

func (r *articleRepository) FindByID(id uint) (*domain.Article, error) {
	var article domain.Article
	if err := r.db.First(&article, id).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) FindAll(page, limit int, search string) ([]domain.Article, int64, error) {
	var articles []domain.Article
	var total int64

	query := r.db.Model(&domain.Article{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("title ILIKE ? OR slug ILIKE ? OR content ILIKE ?", like, like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

func (r *articleRepository) Update(article *domain.Article) error {
	return r.db.Save(article).Error
}

func (r *articleRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Article{}, id).Error
}

