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

// MockCommentService is a mock implementation of service.CommentService
type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) CreateComment(userID, bookID uint, req dto.CreateCommentRequest) (*domain.Comment, error) {
	args := m.Called(userID, bookID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockCommentService) ListComments(bookID uint) ([]domain.Comment, error) {
	args := m.Called(bookID)
	return args.Get(0).([]domain.Comment), args.Error(1)
}

func (m *MockCommentService) UpdateComment(userID, commentID uint, userRole domain.UserRole, req dto.UpdateCommentRequest) (*domain.Comment, error) {
	args := m.Called(userID, commentID, userRole, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockCommentService) DeleteComment(userID, commentID uint, userRole domain.UserRole) error {
	args := m.Called(userID, commentID, userRole)
	return args.Error(0)
}

func TestCreateComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockCommentService)
	handler := NewCommentHandler(mockService)

	router := gin.Default()
	router.POST("/books/:id/comments", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		handler.CreateComment(c)
	})

	reqBody := dto.CreateCommentRequest{Content: "Great"}
	jsonBody, _ := json.Marshal(reqBody)

	expectedComment := &domain.Comment{ID: 1, Content: "Great"}

	mockService.On("CreateComment", uint(1), uint(1), reqBody).Return(expectedComment, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/books/1/comments", bytes.NewBuffer(jsonBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockService.AssertExpectations(t)
}

func TestListComments(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockCommentService)
	handler := NewCommentHandler(mockService)

	router := gin.Default()
	router.GET("/books/:id/comments", handler.ListComments)

	expectedComments := []domain.Comment{{ID: 1, Content: "Great"}}

	mockService.On("ListComments", uint(1)).Return(expectedComments, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/1/comments", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockCommentService)
	handler := NewCommentHandler(mockService)

	router := gin.Default()
	router.PUT("/books/:id/comments/:comment_id", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("role", "STUDENT")
		handler.UpdateComment(c)
	})

	reqBody := dto.UpdateCommentRequest{Content: "Better"}
	jsonBody, _ := json.Marshal(reqBody)

	expectedComment := &domain.Comment{ID: 1, Content: "Better"}

	mockService.On("UpdateComment", uint(1), uint(1), domain.RoleStudent, reqBody).Return(expectedComment, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/books/1/comments/1", bytes.NewBuffer(jsonBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestUpdateComment_Forbidden(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockCommentService)
	handler := NewCommentHandler(mockService)

	router := gin.Default()
	router.PUT("/books/:id/comments/:comment_id", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("role", "STUDENT")
		handler.UpdateComment(c)
	})

	reqBody := dto.UpdateCommentRequest{Content: "Better"}
	jsonBody, _ := json.Marshal(reqBody)

	mockService.On("UpdateComment", uint(1), uint(1), domain.RoleStudent, reqBody).Return(nil, errors.New("unauthorized to update this comment"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/books/1/comments/1", bytes.NewBuffer(jsonBody))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
	mockService.AssertExpectations(t)
}

func TestDeleteComment(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockCommentService)
	handler := NewCommentHandler(mockService)

	router := gin.Default()
	router.DELETE("/books/:id/comments/:comment_id", func(c *gin.Context) {
		c.Set("user_id", uint(1))
		c.Set("role", "ADMIN")
		handler.DeleteComment(c)
	})

	mockService.On("DeleteComment", uint(1), uint(1), domain.RoleAdmin).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/books/1/comments/1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}
