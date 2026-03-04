package service

import (
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockReviewRepository is a mock implementation of repository.ReviewRepository
type MockReviewRepository struct {
	mock.Mock
}

func (m *MockReviewRepository) Create(review *domain.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockReviewRepository) FindByID(id uint) (*domain.Review, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Review), args.Error(1)
}

func (m *MockReviewRepository) FindByBookID(bookID uint) ([]domain.Review, error) {
	args := m.Called(bookID)
	return args.Get(0).([]domain.Review), args.Error(1)
}

func (m *MockReviewRepository) Update(review *domain.Review) error {
	args := m.Called(review)
	return args.Error(0)
}

func (m *MockReviewRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestCreateReview(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewService(mockRepo)

	req := dto.CreateReviewRequest{
		Rating:  5,
		Comment: "Great book!",
	}

	mockRepo.On("Create", mock.AnythingOfType("*domain.Review")).Return(nil)

	review, err := service.CreateReview(1, 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, review)
	assert.Equal(t, req.Rating, review.Rating)
	assert.Equal(t, req.Comment, review.Comment)
	mockRepo.AssertExpectations(t)
}

func TestListReviews(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewService(mockRepo)

	expectedReviews := []domain.Review{{ID: 1, Comment: "Good"}}

	mockRepo.On("FindByBookID", uint(1)).Return(expectedReviews, nil)

	reviews, err := service.ListReviews(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedReviews, reviews)
	mockRepo.AssertExpectations(t)
}

func TestUpdateReview_Owner(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewService(mockRepo)

	existingReview := &domain.Review{ID: 1, UserID: 1, Rating: 4, Comment: "Good"}
	req := dto.UpdateReviewRequest{Rating: 5}

	mockRepo.On("FindByID", uint(1)).Return(existingReview, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Review")).Return(nil)

	review, err := service.UpdateReview(1, 1, domain.RoleStudent, req)

	assert.NoError(t, err)
	assert.Equal(t, 5, review.Rating)
	mockRepo.AssertExpectations(t)
}

func TestUpdateReview_Unauthorized(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewService(mockRepo)

	existingReview := &domain.Review{ID: 1, UserID: 2, Rating: 4, Comment: "Good"}
	req := dto.UpdateReviewRequest{Rating: 5}

	mockRepo.On("FindByID", uint(1)).Return(existingReview, nil)

	review, err := service.UpdateReview(1, 1, domain.RoleStudent, req)

	assert.Error(t, err)
	assert.Nil(t, review)
	assert.Equal(t, "unauthorized to update this review", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeleteReview_Admin(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewService(mockRepo)

	existingReview := &domain.Review{ID: 1, UserID: 2}

	mockRepo.On("FindByID", uint(1)).Return(existingReview, nil)
	mockRepo.On("Delete", uint(1)).Return(nil)

	err := service.DeleteReview(1, 1, domain.RoleAdmin)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteReview_Unauthorized(t *testing.T) {
	mockRepo := new(MockReviewRepository)
	service := NewReviewService(mockRepo)

	existingReview := &domain.Review{ID: 1, UserID: 2}

	mockRepo.On("FindByID", uint(1)).Return(existingReview, nil)

	err := service.DeleteReview(1, 1, domain.RoleStudent)

	assert.Error(t, err)
	assert.Equal(t, "unauthorized to delete this review", err.Error())
	mockRepo.AssertExpectations(t)
}
