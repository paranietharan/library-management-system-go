package repository

import (
	"library-management-system-go/internal/domain"

	"gorm.io/gorm"
)

type FineRepository interface {
	Create(fine *domain.Fine) error
	FindByID(id uint) (*domain.Fine, error)
	FindAll(page, limit int, search string) ([]domain.Fine, int64, error)
}

type fineRepository struct {
	db *gorm.DB
}

func NewFineRepository(db *gorm.DB) FineRepository {
	return &fineRepository{db: db}
}

func (r *fineRepository) Create(fine *domain.Fine) error {
	return r.db.Create(fine).Error
}

func (r *fineRepository) FindByID(id uint) (*domain.Fine, error) {
	var fine domain.Fine
	if err := r.db.First(&fine, id).Error; err != nil {
		return nil, err
	}
	return &fine, nil
}

func (r *fineRepository) FindAll(page, limit int, search string) ([]domain.Fine, int64, error) {
	var fines []domain.Fine
	var total int64

	query := r.db.Model(&domain.Fine{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("reason ILIKE ? OR status ILIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&fines).Error; err != nil {
		return nil, 0, err
	}
	return fines, total, nil
}

