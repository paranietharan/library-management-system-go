-- Rollback: remove article sub-resources and book review/comment tables.

DROP TRIGGER IF EXISTS update_article_ratings_updated_at ON article_ratings;
DROP TABLE IF EXISTS article_ratings;

DROP TRIGGER IF EXISTS update_article_comments_updated_at ON article_comments;
DROP TABLE IF EXISTS article_comments;

DROP TRIGGER IF EXISTS update_article_reviews_updated_at ON article_reviews;
DROP TABLE IF EXISTS article_reviews;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'article_review_status') THEN
        DROP TYPE article_review_status;
    END IF;
END
$$;

DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP TABLE IF EXISTS comments;

DROP TRIGGER IF EXISTS update_reviews_updated_at ON reviews;
DROP TABLE IF EXISTS reviews;

