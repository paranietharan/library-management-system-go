package handler

import (
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	service service.ReviewService
}

func NewReviewHandler(service service.ReviewService) *ReviewHandler {
	return &ReviewHandler{service: service}
}

// @Summary Create book review
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param request body dto.CreateReviewRequest true "Create review request"
// @Success 201 {object} gin.H
// @Router /books/{id}/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	userID := c.GetUint("user_id")

	var req dto.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.service.CreateReview(userID, uint(bookID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// @Summary List book reviews
// @Tags reviews
// @Security BearerAuth
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} gin.H
// @Router /books/{id}/reviews [get]
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	reviews, err := h.service.ListReviews(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// @Summary Update book review
// @Tags reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param review_id path int true "Review ID"
// @Param request body dto.UpdateReviewRequest true "Update review request"
// @Success 200 {object} gin.H
// @Router /books/{id}/reviews/{review_id} [put]
func (h *ReviewHandler) UpdateReview(c *gin.Context) {
	reviewID, err := strconv.ParseUint(c.Param("review_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID := c.GetUint("user_id")
	userRole := domain.UserRole(c.GetString("role"))

	var req dto.UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.service.UpdateReview(userID, uint(reviewID), userRole, req)
	if err != nil {
		if err.Error() == "unauthorized to update this review" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// @Summary Delete book review
// @Tags reviews
// @Security BearerAuth
// @Produce json
// @Param id path int true "Book ID"
// @Param review_id path int true "Review ID"
// @Success 200 {object} gin.H
// @Router /books/{id}/reviews/{review_id} [delete]
func (h *ReviewHandler) DeleteReview(c *gin.Context) {
	reviewID, err := strconv.ParseUint(c.Param("review_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID := c.GetUint("user_id")
	userRole := domain.UserRole(c.GetString("role"))

	if err := h.service.DeleteReview(userID, uint(reviewID), userRole); err != nil {
		if err.Error() == "unauthorized to delete this review" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
