package service

import (
	"errors"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBookRepository is a mock implementation of repository.BookRepository
type MockBookRepository struct {
	mock.Mock
}

func (m *MockBookRepository) Create(book *domain.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) FindByID(id uint) (*domain.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (m *MockBookRepository) FindAll(page, limit int, search string) ([]domain.Book, int64, error) {
	args := m.Called(page, limit, search)
	return args.Get(0).([]domain.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookRepository) Update(book *domain.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateBook(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	req := dto.CreateBookRequest{
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "1234567890",
		PublicationYear: 2023,
		Category:        "Fiction",
		TotalCopies:     5,
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Book")).Return(nil)

	book, err := service.CreateBook(req)

	assert.NoError(t, err)
	assert.NotNil(t, book)
	assert.Equal(t, req.Title, book.Title)
	assert.Equal(t, req.Author, book.Author)
	assert.Equal(t, domain.BookStatusAvailable, book.Status)
	mockRepo.AssertExpectations(t)
}

func TestGetBook(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	expectedBook := &domain.Book{ID: 1, Title: "Test Book"}

	mockRepo.On("FindByID", uint(1)).Return(expectedBook, nil)

	book, err := service.GetBook(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedBook, book)
	mockRepo.AssertExpectations(t)
}

func TestGetBook_NotFound(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	mockRepo.On("FindByID", uint(1)).Return(nil, errors.New("not found"))

	book, err := service.GetBook(1)

	assert.Error(t, err)
	assert.Nil(t, book)
	mockRepo.AssertExpectations(t)
}

func TestListBooks(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	expectedBooks := []domain.Book{{ID: 1, Title: "Book 1"}, {ID: 2, Title: "Book 2"}}
	expectedTotal := int64(2)

	mockRepo.On("FindAll", 1, 10, "").Return(expectedBooks, expectedTotal, nil)

	books, total, err := service.ListBooks(1, 10, "")

	assert.NoError(t, err)
	assert.Equal(t, expectedBooks, books)
	assert.Equal(t, expectedTotal, total)
	mockRepo.AssertExpectations(t)
}

func TestUpdateBook(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	existingBook := &domain.Book{ID: 1, Title: "Old Title", TotalCopies: 5, AvailableCopies: 5}
	req := dto.UpdateBookRequest{Title: "New Title"}

	mockRepo.On("FindByID", uint(1)).Return(existingBook, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Book")).Return(nil)

	book, err := service.UpdateBook(1, req)

	assert.NoError(t, err)
	assert.Equal(t, "New Title", book.Title)
	mockRepo.AssertExpectations(t)
}

func TestDeleteBook(t *testing.T) {
	mockRepo := new(MockBookRepository)
	service := NewBookService(mockRepo)

	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.DeleteBook(1)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
