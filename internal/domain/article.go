package domain

import (
	"time"

	"github.com/lib/pq"
)

type Article struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	Title             string    `gorm:"not null" json:"title"`
	Slug              string    `gorm:"unique;not null" json:"slug"`
	Category          string    `gorm:"type:article_category;not null" json:"category"`
	Content           string    `gorm:"not null" json:"content"`
	Excerpt           *string   `json:"excerpt,omitempty"`
	AuthorID          uint      `gorm:"not null" json:"author_id"`
	FeaturedImageURL  *string   `json:"featured_image_url,omitempty"`
	IsPublished       bool      `gorm:"default:false" json:"is_published"`
	PublishedAt       *time.Time `json:"published_at,omitempty"`
	ViewsCount        int       `gorm:"default:0" json:"views_count"`
	Tags              pq.StringArray `gorm:"type:text[]" json:"tags,omitempty"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ArticleReviewStatus string

const (
	ArticleReviewStatusPending  ArticleReviewStatus = "PENDING"
	ArticleReviewStatusApproved ArticleReviewStatus = "APPROVED"
	ArticleReviewStatusRejected ArticleReviewStatus = "REJECTED"
)

type ArticleReview struct {
	ID        uint                   `gorm:"primaryKey" json:"id"`
	ArticleID uint                   `gorm:"not null" json:"article_id"`
	UserID    uint                   `gorm:"not null" json:"user_id"`
	Status    ArticleReviewStatus  `gorm:"type:article_review_status;not null;default:'PENDING'" json:"status"`
	Feedback  *string                `json:"feedback,omitempty"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

type ArticleComment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"not null" json:"article_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Content   string    `gorm:"not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ArticleRating struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ArticleID uint      `gorm:"not null" json:"article_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Rating    int       `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

