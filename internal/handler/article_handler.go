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

// @Summary List articles
// @Tags articles
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param search query string false "Search"
// @Success 200 {object} gin.H
// @Router /articles [get]
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

// @Summary Get article
// @Tags articles
// @Security BearerAuth
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} gin.H
// @Router /articles/{id} [get]
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

// @Summary Create article
// @Tags articles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateArticleRequest true "Create article request"
// @Success 201 {object} gin.H
// @Router /articles [post]
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

// @Summary Update article
// @Tags articles
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param request body dto.UpdateArticleRequest true "Update article request"
// @Success 200 {object} gin.H
// @Router /articles/{id} [put]
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

// @Summary Delete article
// @Tags articles
// @Security BearerAuth
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} gin.H
// @Router /articles/{id} [delete]
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

// @Summary List article reviews
// @Tags article-reviews
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page"
// @Param limit query int false "Limit"
// @Param search query string false "Search"
// @Success 200 {object} gin.H
// @Router /articles/review [get]
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

// @Summary Get article review
// @Tags article-reviews
// @Security BearerAuth
// @Produce json
// @Param id path int true "Review ID"
// @Success 200 {object} gin.H
// @Router /articles/review/{id} [get]
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

// @Summary Create article review
// @Tags article-reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateArticleReviewRequest true "Create article review request"
// @Success 201 {object} gin.H
// @Router /articles/review [post]
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

// @Summary Update article review
// @Tags article-reviews
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Review ID"
// @Param request body dto.UpdateArticleReviewRequest true "Update article review request"
// @Success 200 {object} gin.H
// @Router /articles/review/{id} [put]
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

// @Summary List article comments
// @Tags article-comments
// @Security BearerAuth
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} gin.H
// @Router /articles/{id}/comments [get]
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

// @Summary Create article comment
// @Tags article-comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param request body dto.CreateArticleCommentRequest true "Create article comment request"
// @Success 201 {object} gin.H
// @Router /articles/{id}/comments [post]
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

// @Summary Update article comment
// @Tags article-comments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param comment_id path int true "Comment ID"
// @Param request body dto.UpdateArticleCommentRequest true "Update article comment request"
// @Success 200 {object} gin.H
// @Router /articles/{id}/comments/{comment_id} [put]
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

// @Summary Delete article comment
// @Tags article-comments
// @Security BearerAuth
// @Produce json
// @Param id path int true "Article ID"
// @Param comment_id path int true "Comment ID"
// @Success 200 {object} gin.H
// @Router /articles/{id}/comments/{comment_id} [delete]
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

// @Summary List article ratings
// @Tags article-ratings
// @Security BearerAuth
// @Produce json
// @Param id path int true "Article ID"
// @Success 200 {object} gin.H
// @Router /articles/{id}/ratings [get]
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

// @Summary Create article rating
// @Tags article-ratings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param request body dto.CreateArticleRatingRequest true "Create article rating request"
// @Success 201 {object} gin.H
// @Router /articles/{id}/ratings [post]
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

// @Summary Update article rating
// @Tags article-ratings
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param rating_id path int true "Rating ID"
// @Param request body dto.UpdateArticleRatingRequest true "Update article rating request"
// @Success 200 {object} gin.H
// @Router /articles/{id}/ratings/{rating_id} [put]
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

// @Summary Delete article rating
// @Tags article-ratings
// @Security BearerAuth
// @Produce json
// @Param id path int true "Article ID"
// @Param rating_id path int true "Rating ID"
// @Success 200 {object} gin.H
// @Router /articles/{id}/ratings/{rating_id} [delete]
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
