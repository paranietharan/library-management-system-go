package service

import (
	"errors"
	"fmt"
	"library-management-system-go/internal/domain"
	"library-management-system-go/internal/dto"
	"library-management-system-go/internal/repository"
	"strings"
	"time"
)

type ArticleService interface {
	CreateArticle(userID uint, role domain.UserRole, req dto.CreateArticleRequest) (*domain.Article, error)
	GetArticle(id uint) (*domain.Article, error)
	ListArticles(page, limit int, search string) ([]domain.Article, int64, error)
	UpdateArticle(userID uint, role domain.UserRole, id uint, req dto.UpdateArticleRequest) (*domain.Article, error)
	DeleteArticle(userID uint, role domain.UserRole, id uint) error

	ListArticleReviews(page, limit int, search string) ([]domain.ArticleReview, int64, error)
	GetArticleReview(id uint) (*domain.ArticleReview, error)
	CreateArticleReview(teacherID uint, req dto.CreateArticleReviewRequest) (*domain.ArticleReview, error)
	UpdateArticleReview(teacherID uint, role domain.UserRole, reviewID uint, req dto.UpdateArticleReviewRequest) (*domain.ArticleReview, error)

	ListArticleComments(articleID uint) ([]domain.ArticleComment, error)
	CreateArticleComment(userID uint, role domain.UserRole, articleID uint, req dto.CreateArticleCommentRequest) (*domain.ArticleComment, error)
	UpdateArticleComment(userID uint, role domain.UserRole, commentID uint, req dto.UpdateArticleCommentRequest) (*domain.ArticleComment, error)
	DeleteArticleComment(userID uint, role domain.UserRole, commentID uint) error

	ListArticleRatings(articleID uint) ([]domain.ArticleRating, error)
	CreateArticleRating(userID uint, role domain.UserRole, articleID uint, req dto.CreateArticleRatingRequest) (*domain.ArticleRating, error)
	UpdateArticleRating(userID uint, role domain.UserRole, ratingID uint, req dto.UpdateArticleRatingRequest) (*domain.ArticleRating, error)
	DeleteArticleRating(userID uint, role domain.UserRole, ratingID uint) error
}

type articleService struct {
	articleRepo        repository.ArticleRepository
	articleReviewRepo repository.ArticleReviewRepository
	commentRepo        repository.ArticleCommentRepository
	ratingRepo         repository.ArticleRatingRepository
}

func NewArticleService(
	articleRepo repository.ArticleRepository,
	articleReviewRepo repository.ArticleReviewRepository,
	commentRepo repository.ArticleCommentRepository,
	ratingRepo repository.ArticleRatingRepository,
) ArticleService {
	return &articleService{
		articleRepo:        articleRepo,
		articleReviewRepo: articleReviewRepo,
		commentRepo:        commentRepo,
		ratingRepo:         ratingRepo,
	}
}

func (s *articleService) CreateArticle(userID uint, role domain.UserRole, req dto.CreateArticleRequest) (*domain.Article, error) {
	if role != domain.RoleStudent && role != domain.RoleAdmin {
		return nil, errors.New("only students can submit articles")
	}

	slug := ""
	if req.Slug != nil && strings.TrimSpace(*req.Slug) != "" {
		slug = strings.TrimSpace(*req.Slug)
	} else {
		slug = s.slugify(req.Title)
	}

	tags := req.Tags
	if tags == nil {
		tags = []string{}
	}

	article := &domain.Article{
		Title:            req.Title,
		Slug:             slug,
		Category:        req.Category,
		Content:         req.Content,
		Excerpt:         req.Excerpt,
		AuthorID:        userID,
		FeaturedImageURL: req.FeaturedImageURL,
		IsPublished:     false,
		PublishedAt:     nil,
		ViewsCount:      0,
		Tags:            tags,
	}

	if err := s.articleRepo.Create(article); err != nil {
		return nil, err
	}

	return article, nil
}

func (s *articleService) GetArticle(id uint) (*domain.Article, error) {
	return s.articleRepo.FindByID(id)
}

func (s *articleService) ListArticles(page, limit int, search string) ([]domain.Article, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.articleRepo.FindAll(page, limit, search)
}

func (s *articleService) UpdateArticle(userID uint, role domain.UserRole, id uint, req dto.UpdateArticleRequest) (*domain.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Docs say Student only; teachers/admin can still act via reviews (publish gate).
	if role != domain.RoleStudent && role != domain.RoleAdmin {
		return nil, errors.New("insufficient permissions")
	}
	if role == domain.RoleStudent && article.AuthorID != userID {
		return nil, errors.New("unauthorized to update this article")
	}
	if article.IsPublished {
		return nil, errors.New("published articles cannot be edited")
	}

	if req.Title != nil {
		article.Title = *req.Title
	}
	if req.Slug != nil && strings.TrimSpace(*req.Slug) != "" {
		article.Slug = strings.TrimSpace(*req.Slug)
	}
	if req.Category != nil {
		article.Category = *req.Category
	}
	if req.Content != nil {
		article.Content = *req.Content
	}
	if req.Excerpt != nil {
		article.Excerpt = req.Excerpt
	}
	if req.FeaturedImageURL != nil {
		article.FeaturedImageURL = req.FeaturedImageURL
	}
	if req.Tags != nil {
		article.Tags = *req.Tags
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}
	return article, nil
}

func (s *articleService) DeleteArticle(userID uint, role domain.UserRole, id uint) error {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return err
	}

	switch role {
	case domain.RoleStudent:
		if article.AuthorID != userID {
			return errors.New("unauthorized to delete this article")
		}
		if article.IsPublished {
			return errors.New("published articles cannot be deleted by students")
		}
	case domain.RoleTeacher, domain.RoleAdmin:
		// Allowed
	default:
		return errors.New("insufficient permissions")
	}

	return s.articleRepo.Delete(id)
}

func (s *articleService) ListArticleReviews(page, limit int, search string) ([]domain.ArticleReview, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	return s.articleReviewRepo.FindAll(page, limit, search)
}

func (s *articleService) GetArticleReview(id uint) (*domain.ArticleReview, error) {
	return s.articleReviewRepo.FindByID(id)
}

func (s *articleService) CreateArticleReview(teacherID uint, req dto.CreateArticleReviewRequest) (*domain.ArticleReview, error) {
	status := domain.ArticleReviewStatus(strings.ToUpper(req.Status))
	if status != domain.ArticleReviewStatusPending && status != domain.ArticleReviewStatusApproved && status != domain.ArticleReviewStatusRejected {
		return nil, errors.New("invalid review status")
	}

	review := &domain.ArticleReview{
		ArticleID: req.ArticleID,
		UserID:    teacherID,
		Status:    status,
		Feedback:  req.Feedback,
	}

	if err := s.articleReviewRepo.Create(review); err != nil {
		return nil, err
	}

	// Publish gate
	article, err := s.articleRepo.FindByID(req.ArticleID)
	if err != nil {
		return nil, err
	}

	article.IsPublished = status == domain.ArticleReviewStatusApproved
	if article.IsPublished {
		now := time.Now()
		article.PublishedAt = &now
	} else {
		article.PublishedAt = nil
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *articleService) UpdateArticleReview(teacherID uint, role domain.UserRole, reviewID uint, req dto.UpdateArticleReviewRequest) (*domain.ArticleReview, error) {
	review, err := s.articleReviewRepo.FindByID(reviewID)
	if err != nil {
		return nil, err
	}

	if role != domain.RoleAdmin && review.UserID != teacherID {
		return nil, errors.New("unauthorized to update this review")
	}

	var newStatus domain.ArticleReviewStatus = review.Status
	if req.Status != nil && strings.TrimSpace(*req.Status) != "" {
		newStatus = domain.ArticleReviewStatus(strings.ToUpper(*req.Status))
	}

	if newStatus != domain.ArticleReviewStatusPending && newStatus != domain.ArticleReviewStatusApproved && newStatus != domain.ArticleReviewStatusRejected {
		return nil, errors.New("invalid review status")
	}

	review.Status = newStatus
	if req.Feedback != nil {
		review.Feedback = req.Feedback
	}

	if err := s.articleReviewRepo.Update(review); err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	article, err := s.articleRepo.FindByID(review.ArticleID)
	if err != nil {
		return nil, err
	}
	article.IsPublished = review.Status == domain.ArticleReviewStatusApproved
	if article.IsPublished {
		now := time.Now()
		article.PublishedAt = &now
	} else {
		article.PublishedAt = nil
	}

	if err := s.articleRepo.Update(article); err != nil {
		return nil, err
	}

	return review, nil
}

func (s *articleService) ListArticleComments(articleID uint) ([]domain.ArticleComment, error) {
	return s.commentRepo.FindByArticleID(articleID)
}

func (s *articleService) CreateArticleComment(userID uint, role domain.UserRole, articleID uint, req dto.CreateArticleCommentRequest) (*domain.ArticleComment, error) {
	switch role {
	case domain.RoleStudent, domain.RoleTeacher, domain.RoleAdmin:
		// allowed
	default:
		return nil, errors.New("insufficient permissions to comment")
	}

	comment := &domain.ArticleComment{
		ArticleID: articleID,
		UserID:    userID,
		Content:   req.Content,
	}

	if err := s.commentRepo.Create(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *articleService) UpdateArticleComment(userID uint, role domain.UserRole, commentID uint, req dto.UpdateArticleCommentRequest) (*domain.ArticleComment, error) {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return nil, err
	}

	if comment.UserID != userID && role != domain.RoleAdmin && role != domain.RoleLibrarian {
		return nil, errors.New("unauthorized to update this comment")
	}

	comment.Content = req.Content
	if err := s.commentRepo.Update(comment); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *articleService) DeleteArticleComment(userID uint, role domain.UserRole, commentID uint) error {
	comment, err := s.commentRepo.FindByID(commentID)
	if err != nil {
		return err
	}

	if comment.UserID != userID && role != domain.RoleAdmin && role != domain.RoleLibrarian {
		return errors.New("unauthorized to delete this comment")
	}

	return s.commentRepo.Delete(commentID)
}

func (s *articleService) ListArticleRatings(articleID uint) ([]domain.ArticleRating, error) {
	return s.ratingRepo.FindByArticleID(articleID)
}

func (s *articleService) CreateArticleRating(userID uint, role domain.UserRole, articleID uint, req dto.CreateArticleRatingRequest) (*domain.ArticleRating, error) {
	switch role {
	case domain.RoleStudent, domain.RoleTeacher, domain.RoleAdmin:
	default:
		return nil, errors.New("insufficient permissions to rate")
	}

	rating := &domain.ArticleRating{
		ArticleID: articleID,
		UserID:    userID,
		Rating:    req.Rating,
	}

	if err := s.ratingRepo.Create(rating); err != nil {
		return nil, err
	}
	return rating, nil
}

func (s *articleService) UpdateArticleRating(userID uint, role domain.UserRole, ratingID uint, req dto.UpdateArticleRatingRequest) (*domain.ArticleRating, error) {
	rating, err := s.ratingRepo.FindByID(ratingID)
	if err != nil {
		return nil, err
	}

	if rating.UserID != userID && role != domain.RoleAdmin && role != domain.RoleLibrarian {
		return nil, errors.New("unauthorized to update this rating")
	}

	rating.Rating = req.Rating
	if err := s.ratingRepo.Update(rating); err != nil {
		return nil, err
	}
	return rating, nil
}

func (s *articleService) DeleteArticleRating(userID uint, role domain.UserRole, ratingID uint) error {
	rating, err := s.ratingRepo.FindByID(ratingID)
	if err != nil {
		return err
	}

	if rating.UserID != userID && role != domain.RoleAdmin && role != domain.RoleLibrarian {
		return errors.New("unauthorized to delete this rating")
	}

	return s.ratingRepo.Delete(ratingID)
}

func (s *articleService) slugify(s0 string) string {
	s0 = strings.ToLower(strings.TrimSpace(s0))
	// naive slugification; good enough for now.
	s0 = strings.ReplaceAll(s0, " ", "-")
	s0 = strings.ReplaceAll(s0, "_", "-")
	return s0
}

