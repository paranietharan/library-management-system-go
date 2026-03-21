package repository

import (
	"library-management-system-go/internal/domain"

	"gorm.io/gorm"
)

type ArticleCommentRepository interface {
	Create(comment *domain.ArticleComment) error
	FindByID(id uint) (*domain.ArticleComment, error)
	FindByArticleID(articleID uint) ([]domain.ArticleComment, error)
	Update(comment *domain.ArticleComment) error
	Delete(id uint) error
}

type articleCommentRepository struct {
	db *gorm.DB
}

func NewArticleCommentRepository(db *gorm.DB) ArticleCommentRepository {
	return &articleCommentRepository{db: db}
}

func (r *articleCommentRepository) Create(comment *domain.ArticleComment) error {
	return r.db.Create(comment).Error
}

func (r *articleCommentRepository) FindByID(id uint) (*domain.ArticleComment, error) {
	var comment domain.ArticleComment
	if err := r.db.First(&comment, id).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *articleCommentRepository) FindByArticleID(articleID uint) ([]domain.ArticleComment, error) {
	var comments []domain.ArticleComment
	if err := r.db.Where("article_id = ?", articleID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *articleCommentRepository) Update(comment *domain.ArticleComment) error {
	return r.db.Save(comment).Error
}

func (r *articleCommentRepository) Delete(id uint) error {
	return r.db.Delete(&domain.ArticleComment{}, id).Error
}

