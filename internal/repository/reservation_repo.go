package repository

import (
	"library-management-system-go/internal/domain"

	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(reservation *domain.Reservation) error
	FindByID(id uint) (*domain.Reservation, error)
	FindAll(page, limit int, search string) ([]domain.Reservation, int64, error)
	Update(reservation *domain.Reservation) error
	Delete(id uint) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(reservation *domain.Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *reservationRepository) FindByID(id uint) (*domain.Reservation, error) {
	var reservation domain.Reservation
	if err := r.db.First(&reservation, id).Error; err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) FindAll(page, limit int, search string) ([]domain.Reservation, int64, error) {
	var reservations []domain.Reservation
	var total int64

	query := r.db.Model(&domain.Reservation{})
	if search != "" {
		like := "%" + search + "%"
		query = query.Joins("JOIN books ON books.id = reservations.book_id").
			Where("books.title ILIKE ? OR books.author ILIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&reservations).Error; err != nil {
		return nil, 0, err
	}
	return reservations, total, nil
}

func (r *reservationRepository) Update(reservation *domain.Reservation) error {
	return r.db.Save(reservation).Error
}

func (r *reservationRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Reservation{}, id).Error
}

