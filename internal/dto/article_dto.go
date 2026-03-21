package dto

import "time"

type CreateArticleRequest struct {
	Title            string   `json:"title" binding:"required,min=3,max=255"`
	Slug             *string  `json:"slug,omitempty"`
	Category         string   `json:"category" binding:"required"`
	Content          string   `json:"content" binding:"required"`
	Excerpt          *string  `json:"excerpt,omitempty"`
	FeaturedImageURL *string  `json:"featured_image_url,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

type UpdateArticleRequest struct {
	Title            *string  `json:"title,omitempty"`
	Slug             *string  `json:"slug,omitempty"`
	Category         *string  `json:"category,omitempty"`
	Content          *string  `json:"content,omitempty"`
	Excerpt          *string  `json:"excerpt,omitempty"`
	FeaturedImageURL *string  `json:"featured_image_url,omitempty"`
	Tags             *[]string `json:"tags,omitempty"`
}

type ArticleResponse struct {
	ID               uint      `json:"id"`
	Title            string    `json:"title"`
	Slug             string    `json:"slug"`
	Category         string    `json:"category"`
	Content          string    `json:"content"`
	Excerpt          *string   `json:"excerpt,omitempty"`
	AuthorID         uint      `json:"author_id"`
	FeaturedImageURL *string   `json:"featured_image_url,omitempty"`
	IsPublished      bool      `json:"is_published"`
	PublishedAt      *time.Time `json:"published_at,omitempty"`
	ViewsCount       int       `json:"views_count"`
	Tags             []string  `json:"tags,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type CreateArticleReviewRequest struct {
	ArticleID uint   `json:"article_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
	Feedback  *string `json:"feedback,omitempty"`
}

type UpdateArticleReviewRequest struct {
	Status   *string `json:"status,omitempty"`
	Feedback *string `json:"feedback,omitempty"`
}

type ArticleReviewResponse struct {
	ID        uint      `json:"id"`
	ArticleID uint      `json:"article_id"`
	UserID    uint      `json:"user_id"`
	Status    string    `json:"status"`
	Feedback  *string   `json:"feedback,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateArticleCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type UpdateArticleCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

type ArticleCommentResponse struct {
	ID        uint      `json:"id"`
	ArticleID uint      `json:"article_id"`
	UserID    uint      `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateArticleRatingRequest struct {
	Rating int `json:"rating" binding:"required,min=1,max=5"`
}

type UpdateArticleRatingRequest struct {
	Rating int `json:"rating" binding:"required,min=1,max=5"`
}

type ArticleRatingResponse struct {
	ID        uint      `json:"id"`
	ArticleID uint      `json:"article_id"`
	UserID    uint      `json:"user_id"`
	Rating    int       `json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

