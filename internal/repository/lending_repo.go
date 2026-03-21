package repository

import (
	"library-management-system-go/internal/domain"

	"gorm.io/gorm"
)

type LendingRepository interface {
	Create(lending *domain.Lending) error
	FindByID(id uint) (*domain.Lending, error)
	FindAll(page, limit int, search string) ([]domain.Lending, int64, error)
	Update(lending *domain.Lending) error
	Delete(id uint) error
}

type lendingRepository struct {
	db *gorm.DB
}

func NewLendingRepository(db *gorm.DB) LendingRepository {
	return &lendingRepository{db: db}
}

func (r *lendingRepository) Create(lending *domain.Lending) error {
	return r.db.Create(lending).Error
}

func (r *lendingRepository) FindByID(id uint) (*domain.Lending, error) {
	var lending domain.Lending
	if err := r.db.First(&lending, id).Error; err != nil {
		return nil, err
	}
	return &lending, nil
}

func (r *lendingRepository) FindAll(page, limit int, search string) ([]domain.Lending, int64, error) {
	var lendings []domain.Lending
	var total int64

	query := r.db.Model(&domain.Lending{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Joins("JOIN books ON books.id = lendings.book_id").
			Where("books.title ILIKE ? OR books.author ILIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&lendings).Error; err != nil {
		return nil, 0, err
	}
	return lendings, total, nil
}

func (r *lendingRepository) Update(lending *domain.Lending) error {
	return r.db.Save(lending).Error
}

func (r *lendingRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Lending{}, id).Error
}

