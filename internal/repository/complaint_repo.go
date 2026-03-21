package repository

import (
	"library-management-system-go/internal/domain"

	"gorm.io/gorm"
)

type ComplaintRepository interface {
	Create(complaint *domain.Complaint) error
	FindByID(id uint) (*domain.Complaint, error)
	FindAll(page, limit int, search string) ([]domain.Complaint, int64, error)
	Update(complaint *domain.Complaint) error
	Delete(id uint) error
}

type complaintRepository struct {
	db *gorm.DB
}

func NewComplaintRepository(db *gorm.DB) ComplaintRepository {
	return &complaintRepository{db: db}
}

func (r *complaintRepository) Create(complaint *domain.Complaint) error {
	return r.db.Create(complaint).Error
}

func (r *complaintRepository) FindByID(id uint) (*domain.Complaint, error) {
	var complaint domain.Complaint
	if err := r.db.First(&complaint, id).Error; err != nil {
		return nil, err
	}
	return &complaint, nil
}

func (r *complaintRepository) FindAll(page, limit int, search string) ([]domain.Complaint, int64, error) {
	var complaints []domain.Complaint
	var total int64

	query := r.db.Model(&domain.Complaint{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Where("subject ILIKE ? OR description ILIKE ? OR status ILIKE ?", like, like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&complaints).Error; err != nil {
		return nil, 0, err
	}

	return complaints, total, nil
}

func (r *complaintRepository) Update(complaint *domain.Complaint) error {
	return r.db.Save(complaint).Error
}

func (r *complaintRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Complaint{}, id).Error
}

