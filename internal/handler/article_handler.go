package handler

import (
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	service service.ArticleService
}

func NewArticleHandler(service service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service: service}
}

func (h *ArticleHandler) ListArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	articles, total, err := h.service.ListArticles(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  articles,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *ArticleHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	article, err := h.service.GetArticle(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.CreateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article, err := h.service.CreateArticle(userID, role, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.UpdateArticleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article, err := h.service.UpdateArticle(userID, role, uint(id), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	if err := h.service.DeleteArticle(userID, role, uint(id)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func (h *ArticleHandler) ListArticleReviews(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	reviews, total, err := h.service.ListArticleReviews(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  reviews,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *ArticleHandler) GetArticleReview(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	review, err := h.service.GetArticleReview(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	c.JSON(http.StatusOK, review)
}

func (h *ArticleHandler) CreateArticleReview(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req dto.CreateArticleReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.service.CreateArticleReview(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *ArticleHandler) UpdateArticleReview(c *gin.Context) {
	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.UpdateArticleReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.service.UpdateArticleReview(userID, role, uint(reviewID), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

func (h *ArticleHandler) ListArticleComments(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	comments, err := h.service.ListArticleComments(uint(articleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *ArticleHandler) CreateArticleComment(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.CreateArticleCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	comment, err := h.service.CreateArticleComment(userID, role, uint(articleID), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *ArticleHandler) UpdateArticleComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.UpdateArticleCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_ = articleID // article_id is available for validation/extension; currently only comment id is needed.
	comment, err := h.service.UpdateArticleComment(userID, role, uint(commentID), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *ArticleHandler) DeleteArticleComment(c *gin.Context) {
	commentID, err := strconv.ParseUint(c.Param("comment_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	_ = articleID // article_id is available for validation/extension; currently only comment id is needed.
	if err := h.service.DeleteArticleComment(userID, role, uint(commentID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article comment deleted successfully"})
}

func (h *ArticleHandler) ListArticleRatings(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	ratings, err := h.service.ListArticleRatings(uint(articleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ratings)
}

func (h *ArticleHandler) CreateArticleRating(c *gin.Context) {
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid article ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.CreateArticleRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rating, err := h.service.CreateArticleRating(userID, role, uint(articleID), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, rating)
}

func (h *ArticleHandler) UpdateArticleRating(c *gin.Context) {
	ratingID, err := strconv.ParseUint(c.Param("rating_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	var req dto.UpdateArticleRatingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rating, err := h.service.UpdateArticleRating(userID, role, uint(ratingID), req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rating)
}

func (h *ArticleHandler) DeleteArticleRating(c *gin.Context) {
	ratingID, err := strconv.ParseUint(c.Param("rating_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rating ID"})
		return
	}

	userID := c.GetUint("user_id")
	role := domain.UserRole(c.GetString("role"))

	if err := h.service.DeleteArticleRating(userID, role, uint(ratingID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article rating deleted successfully"})
}

