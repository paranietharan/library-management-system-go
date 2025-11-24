package service

import (
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCommentRepository is a mock implementation of repository.CommentRepository
type MockCommentRepository struct {
	mock.Mock
}

func (m *MockCommentRepository) Create(comment *domain.Comment) error {
	args := m.Called(comment)
	return args.Error(0)
}

func (m *MockCommentRepository) FindByID(id uint) (*domain.Comment, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Comment), args.Error(1)
}

func (m *MockCommentRepository) FindByBookID(bookID uint) ([]domain.Comment, error) {
	args := m.Called(bookID)
	return args.Get(0).([]domain.Comment), args.Error(1)
}

func (m *MockCommentRepository) Update(comment *domain.Comment) error {
	args := m.Called(comment)
	return args.Error(0)
}

func (m *MockCommentRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateComment(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	service := NewCommentService(mockRepo)

	req := dto.CreateCommentRequest{
		Content: "Great book!",
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Comment")).Return(nil)

	comment, err := service.CreateComment(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, comment)
	assert.Equal(t, req.Content, comment.Content)
	mockRepo.AssertExpectations(t)
}

func TestListComments(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	service := NewCommentService(mockRepo)

	expectedComments := []domain.Comment{{ID: 1, Content: "Good"}}

	mockRepo.On("FindByBookID", uint(1)).Return(expectedComments, nil)

	comments, err := service.ListComments(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedComments, comments)
	mockRepo.AssertExpectations(t)
}

func TestUpdateComment_Owner(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	service := NewCommentService(mockRepo)

	existingComment := &domain.Comment{ID: 1, UserID: 1, Content: "Good"}
	req := dto.UpdateCommentRequest{Content: "Better"}

	mockRepo.On("FindByID", uint(1)).Return(existingComment, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Comment")).Return(nil)

	comment, err := service.UpdateComment(1, 1, domain.RoleStudent, req)

	assert.NoError(t, err)
	assert.Equal(t, "Better", comment.Content)
	mockRepo.AssertExpectations(t)
}

func TestUpdateComment_Unauthorized(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	service := NewCommentService(mockRepo)

	existingComment := &domain.Comment{ID: 1, UserID: 2, Content: "Good"}
	req := dto.UpdateCommentRequest{Content: "Better"}

	mockRepo.On("FindByID", uint(1)).Return(existingComment, nil)

	comment, err := service.UpdateComment(1, 1, domain.RoleStudent, req)

	assert.Error(t, err)
	assert.Nil(t, comment)
	assert.Equal(t, "unauthorized to update this comment", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeleteComment_Admin(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	service := NewCommentService(mockRepo)

	existingComment := &domain.Comment{ID: 1, UserID: 2}

	mockRepo.On("FindByID", uint(1)).Return(existingComment, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.DeleteComment(1, 1, domain.RoleAdmin)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteComment_Unauthorized(t *testing.T) {
	mockRepo := new(MockCommentRepository)
	service := NewCommentService(mockRepo)

	existingComment := &domain.Comment{ID: 1, UserID: 2}

	mockRepo.On("FindByID", uint(1)).Return(existingComment, nil)

	err := service.DeleteComment(1, 1, domain.RoleStudent)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized to delete this comment", err.Error())
	mockRepo.AssertExpectations(t)
}
