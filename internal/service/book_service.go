package service

import (
	"errors"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/repository"
)

type BookService interface {
	CreateBook(req dto.CreateBookRequest) (*domain.Book, error)
	GetBook(id uint) (*domain.Book, error)
	ListBooks(page, limit int, search string) ([]domain.Book, int64, error)
	UpdateBook(id uint, req dto.UpdateBookRequest) (*domain.Book, error)
	DeleteBook(id uint) error
}

type bookService struct {
	repo repository.BookRepository
}

func NewBookService(repo repository.BookRepository) BookService {
	return &bookService{repo: repo}
}

func (s *bookService) CreateBook(req dto.CreateBookRequest) (*domain.Book, error) {
	book := &domain.Book{
		Title:           req.Title,
		Author:          req.Author,
		ISBN:            req.ISBN,
		PublicationYear: req.PublicationYear,
		Category:        req.Category,
		TotalCopies:     req.TotalCopies,
		AvailableCopies: req.TotalCopies,
		Status:          domain.BookStatusAvailable,
	}

	if err := s.repo.Create(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) GetBook(id uint) (*domain.Book, error) {
	return s.repo.FindByID(id)
}

func (s *bookService) ListBooks(page, limit int, search string) ([]domain.Book, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.repo.FindAll(page, limit, search)
}

func (s *bookService) UpdateBook(id uint, req dto.UpdateBookRequest) (*domain.Book, error) {
	book, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.ISBN != "" {
		book.ISBN = req.ISBN
	}
	if req.PublicationYear != 0 {
		book.PublicationYear = req.PublicationYear
	}
	if req.Category != "" {
		book.Category = req.Category
	}
	if req.TotalCopies != 0 {
		// Adjust available copies based on the change in total copies
		diff := req.TotalCopies - book.TotalCopies
		book.TotalCopies = req.TotalCopies
		book.AvailableCopies += diff
		if book.AvailableCopies < 0 {
			return nil, errors.New("cannot reduce total copies below currently borrowed amount")
		}
	}
	if req.Status != "" {
		book.Status = domain.BookStatus(req.Status)
	}

	if err := s.repo.Update(book); err != nil {
		return nil, err
	}

	return book, nil
}

func (s *bookService) DeleteBook(id uint) error {
	return s.repo.Delete(id)
}
