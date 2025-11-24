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

// MockReviewService is a mock implementation of service.ReviewService
type MockReviewService struct {
	mock.Mock
}

func (m *MockReviewService) CreateReview(userID, bookID uint, req dto.CreateReviewRequest) (*domain.Review, error) {
	args := m.Called(userID, bookID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Review), args.Error(1)
}

func (m *MockReviewService) ListReviews(bookID uint) ([]domain.Review, error) {
	args := m.Called(bookID)
	return args.Get(0).([]domain.Review), args.Error(1)
}

func (m *MockReviewService) UpdateReview(userID, reviewID uint, userRole domain.UserRole, req dto.UpdateReviewRequest) (*domain.Review, error) {
	args := m.Called(userID, reviewID, userRole, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Review), args.Error(1)
}

func (m *MockReviewService) DeleteReview(userID, reviewID uint, userRole domain.UserRole) error {
	args := m.Called(userID, reviewID, userRole)
	return args.Error(0)
}

func TestCreateReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockReviewService)
	handler := NewReviewHandler(mockService)

	router := gin.Default()
	router.POST("/books/:book_id/reviews", func(c *gin.Context) {
		c.Set("userID", uint(1))
		handler.CreateReview(c)
	})

	reqBody := dto.CreateReviewRequest{Rating: 5, Comment: "Great"}
	jsonBody, _ := json.Marshal(reqBody)

	expectedReview := &domain.Review{ID: 1, Rating: 5, Comment: "Great"}

	mockService.On("CreateReview", uint(1), uint(1), reqBody).Return(expectedReview, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books/1/reviews", bytes.NewBuffer(jsonBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestListReviews(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockReviewService)
	handler := NewReviewHandler(mockService)

	router := gin.Default()
	router.GET("/books/:book_id/reviews", handler.ListReviews)

	expectedReviews := []domain.Review{{ID: 1, Rating: 5}}

	mockService.On("ListReviews", uint(1)).Return(expectedReviews, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/1/reviews", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockReviewService)
	handler := NewReviewHandler(mockService)

	router := gin.Default()
	router.PUT("/books/:book_id/reviews/:review_id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Set("role", "STUDENT")
		handler.UpdateReview(c)
	})

	reqBody := dto.UpdateReviewRequest{Rating: 4}
	jsonBody, _ := json.Marshal(reqBody)

	expectedReview := &domain.Review{ID: 1, Rating: 4}

	mockService.On("UpdateReview", uint(1), uint(1), domain.RoleStudent, reqBody).Return(expectedReview, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/books/1/reviews/1", bytes.NewBuffer(jsonBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateReview_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockReviewService)
	handler := NewReviewHandler(mockService)

	router := gin.Default()
	router.PUT("/books/:book_id/reviews/:review_id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Set("role", "STUDENT")
		handler.UpdateReview(c)
	})

	reqBody := dto.UpdateReviewRequest{Rating: 4}
	jsonBody, _ := json.Marshal(reqBody)

	mockService.On("UpdateReview", uint(1), uint(1), domain.RoleStudent, reqBody).Return(nil, errors.New("unauthorized to update this review"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/books/1/reviews/1", bytes.NewBuffer(jsonBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteReview(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockReviewService)
	handler := NewReviewHandler(mockService)

	router := gin.Default()
	router.DELETE("/books/:book_id/reviews/:review_id", func(c *gin.Context) {
		c.Set("userID", uint(1))
		c.Set("role", "ADMIN")
		handler.DeleteReview(c)
	})

	mockService.On("DeleteReview", uint(1), uint(1), domain.RoleAdmin).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/books/1/reviews/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
