package repository

import (
	"library-management-system-go/internal/domain"

	"gorm.io/gorm"
)

type BookRepository interface {
	Create(book *domain.Book) error
	FindByID(id uint) (*domain.Book, error)
	FindAll(page, limit int, search string) ([]domain.Book, int64, error)
	Update(book *domain.Book) error
	Delete(id uint) error
}

type bookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) BookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(book *domain.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepository) FindByID(id uint) (*domain.Book, error) {
	var book domain.Book
	err := r.db.First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *bookRepository) FindAll(page, limit int, search string) ([]domain.Book, int64, error) {
	var books []domain.Book
	var total int64

	query := r.db.Model(&domain.Book{})

	if search != "" {
		query = query.Where("title ILIKE ? OR author ILIKE ? OR isbn ILIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = query.Offset(offset).Limit(limit).Find(&books).Error
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}

func (r *bookRepository) Update(book *domain.Book) error {
	return r.db.Save(book).Error
}

func (r *bookRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Book{}, id).Error
}
