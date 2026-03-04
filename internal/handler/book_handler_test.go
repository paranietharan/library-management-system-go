package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBookService is a mock implementation of service.BookService
type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) CreateBook(req dto.CreateBookRequest) (*domain.Book, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (m *MockBookService) GetBook(id uint) (*domain.Book, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (m *MockBookService) ListBooks(page, limit int, search string) ([]domain.Book, int64, error) {
	args := m.Called(page, limit, search)
	return args.Get(0).([]domain.Book), args.Get(1).(int64), args.Error(2)
}

func (m *MockBookService) UpdateBook(id uint, req dto.UpdateBookRequest) (*domain.Book, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Book), args.Error(1)
}

func (m *MockBookService) DeleteBook(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	handler := NewBookHandler(mockService)

	router := gin.Default()
	router.POST("/books", handler.CreateBook)

	reqBody := dto.CreateBookRequest{
		Title:           "Test Book",
		Author:          "Test Author",
		ISBN:            "1234567890",
		PublicationYear: 2023,
		Category:        "Fiction",
		TotalCopies:     5,
	}
	jsonBody, _ := json.Marshal(reqBody)

	expectedBook := &domain.Book{
		ID:          1,
		Title:       reqBody.Title,
		Author:      reqBody.Author,
		TotalCopies: reqBody.TotalCopies,
	}

	mockService.On("CreateBook", reqBody).Return(expectedBook, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(jsonBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetBook(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	handler := NewBookHandler(mockService)

	router := gin.Default()
	router.GET("/books/:id", handler.GetBook)

	expectedBook := &domain.Book{ID: 1, Title: "Test Book"}

	mockService.On("GetBook", uint(1)).Return(expectedBook, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestGetBook_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	handler := NewBookHandler(mockService)

	router := gin.Default()
	router.GET("/books/:id", handler.GetBook)

	mockService.On("GetBook", uint(1)).Return(nil, errors.New("not found"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	mockService.AssertExpectations(t)
}

func TestListBooks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockBookService)
	handler := NewBookHandler(mockService)

	router := gin.Default()
	router.GET("/books", handler.ListBooks)

	expectedBooks := []domain.Book{{ID: 1, Title: "Book 1"}}
	expectedTotal := int64(1)

	mockService.On("ListBooks", 1, 10, "").Return(expectedBooks, expectedTotal, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
