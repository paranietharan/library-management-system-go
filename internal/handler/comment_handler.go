package handler

import (
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	service service.CommentService
}

func NewCommentHandler(service service.CommentService) *CommentHandler {
	return &CommentHandler{service: service}
}

// @Summary Create book comment
// @Tags comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param request body dto.CreateCommentRequest true "Create comment request"
// @Success 201 {object} gin.H
// @Router /books/{id}/comments [post]
func (h *CommentHandler) CreateComment(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	userID := c.GetUint("user_id")

	var req dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.service.CreateComment(userID, uint(bookID), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// @Summary List book comments
// @Tags comments
// @Security BearerAuth
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} gin.H
// @Router /books/{id}/comments [get]
func (h *CommentHandler) ListComments(c *gin.Context) {
	bookID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	comments, err := h.service.ListComments(uint(bookID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// @Summary Update book comment
// @Tags comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param comment_id path int true "Comment ID"
// @Param request body dto.UpdateCommentRequest true "Update comment request"
// @Success 200 {object} gin.H
// @Router /books/{id}/comments/{comment_id} [put]
func (h *CommentHandler) UpdateComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID := c.GetUint("user_id")
	userRole := domain.UserRole(c.GetString("role"))

	var req dto.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.service.UpdateComment(userID, uint(commentID), userRole, req)
	if err != nil {
		if err.Error() == "unauthorized to update this comment" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// @Summary Delete book comment
// @Tags comments
// @Security BearerAuth
// @Produce json
// @Param id path int true "Book ID"
// @Param comment_id path int true "Comment ID"
// @Success 200 {object} gin.H
// @Router /books/{id}/comments/{comment_id} [delete]
func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	userID := c.GetUint("user_id")
	userRole := domain.UserRole(c.GetString("role"))

	if err := h.service.DeleteComment(userID, uint(commentID), userRole); err != nil {
		if err.Error() == "unauthorized to delete this comment" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
