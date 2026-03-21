-- Add missing book review/comment tables and article sub-resources.

-- Book Reviews Table (used by /api/v1/books/:book_id/reviews)
CREATE TABLE IF NOT EXISTS reviews (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE RESTRICT,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_book_id ON reviews(book_id);

CREATE TRIGGER update_reviews_updated_at
BEFORE UPDATE ON reviews
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Book Comments Table (used by /api/v1/books/:book_id/comments)
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE RESTRICT,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_comments_user_id ON comments(user_id);
CREATE INDEX IF NOT EXISTS idx_comments_book_id ON comments(book_id);

CREATE TRIGGER update_comments_updated_at
BEFORE UPDATE ON comments
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Article review decision status
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'article_review_status') THEN
        CREATE TYPE article_review_status AS ENUM ('PENDING', 'APPROVED', 'REJECTED');
    END IF;
END
$$;

-- Article Reviews Table (used by /api/v1/articles/review and publication gating)
CREATE TABLE IF NOT EXISTS article_reviews (
    id SERIAL PRIMARY KEY,
    article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    status article_review_status NOT NULL DEFAULT 'PENDING',
    feedback TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_article_reviews_article_id ON article_reviews(article_id);
CREATE INDEX IF NOT EXISTS idx_article_reviews_user_id ON article_reviews(user_id);

CREATE TRIGGER update_article_reviews_updated_at
BEFORE UPDATE ON article_reviews
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Article Comments Table (used by /api/v1/articles/:article_id/comments)
CREATE TABLE IF NOT EXISTS article_comments (
    id SERIAL PRIMARY KEY,
    article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_article_comments_article_id ON article_comments(article_id);
CREATE INDEX IF NOT EXISTS idx_article_comments_user_id ON article_comments(user_id);

CREATE TRIGGER update_article_comments_updated_at
BEFORE UPDATE ON article_comments
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Article Ratings Table (used by /api/v1/articles/:article_id/ratings)
CREATE TABLE IF NOT EXISTS article_ratings (
    id SERIAL PRIMARY KEY,
    article_id INT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_article_ratings_article_id ON article_ratings(article_id);
CREATE INDEX IF NOT EXISTS idx_article_ratings_user_id ON article_ratings(user_id);

CREATE TRIGGER update_article_ratings_updated_at
BEFORE UPDATE ON article_ratings
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

