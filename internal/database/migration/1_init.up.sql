-- Database: library_management

-- Drop tables if they exist (in reverse order of dependencies)
DROP TABLE IF EXISTS complaints CASCADE;
DROP TABLE IF EXISTS fines CASCADE;
DROP TABLE IF EXISTS lendings CASCADE;
DROP TABLE IF EXISTS reservations CASCADE;
DROP TABLE IF EXISTS articles CASCADE;
DROP TABLE IF EXISTS books CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Create ENUM types
CREATE TYPE user_role AS ENUM ('ADMIN', 'LIBRARIAN', 'STUDENT', 'TEACHER');
CREATE TYPE user_status AS ENUM ('ACTIVE', 'INACTIVE', 'SUSPENDED');
CREATE TYPE book_status AS ENUM ('AVAILABLE', 'BORROWED', 'RESERVED', 'MAINTENANCE', 'LOST');
CREATE TYPE lending_status AS ENUM ('ACTIVE', 'RETURNED', 'OVERDUE', 'LOST');
CREATE TYPE reservation_status AS ENUM ('PENDING', 'FULFILLED', 'CANCELLED', 'EXPIRED');
CREATE TYPE fine_status AS ENUM ('PENDING', 'PAID', 'WAIVED');
CREATE TYPE complaint_status AS ENUM ('OPEN', 'IN_PROGRESS', 'RESOLVED', 'CLOSED');
CREATE TYPE article_category AS ENUM ('NEWS', 'EVENT', 'ANNOUNCEMENT', 'GUIDE');

-- Users Table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    role user_role NOT NULL DEFAULT 'STUDENT',
    status user_status NOT NULL DEFAULT 'ACTIVE',
    phone VARCHAR(20),
    address TEXT,
    date_of_birth DATE,
    student_id VARCHAR(50) UNIQUE,
    employee_id VARCHAR(50) UNIQUE,
    max_books_allowed INT DEFAULT 5,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login TIMESTAMP WITH TIME ZONE,
    CONSTRAINT chk_student_or_employee CHECK (
        (role = 'STUDENT' AND student_id IS NOT NULL) OR
        (role IN ('TEACHER', 'LIBRARIAN', 'ADMIN') AND employee_id IS NOT NULL) OR
        (role = 'STUDENT' AND employee_id IS NULL) OR
        (role IN ('TEACHER', 'LIBRARIAN', 'ADMIN') AND student_id IS NULL)
    )
);

-- Books Table
CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    isbn VARCHAR(13) UNIQUE NOT NULL,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    publisher VARCHAR(255),
    publication_year INT,
    edition VARCHAR(50),
    category VARCHAR(100),
    language VARCHAR(50) DEFAULT 'English',
    pages INT,
    description TEXT,
    location VARCHAR(100),
    shelf_number VARCHAR(50),
    status book_status NOT NULL DEFAULT 'AVAILABLE',
    total_copies INT NOT NULL DEFAULT 1,
    available_copies INT NOT NULL DEFAULT 1,
    cover_image_url TEXT,
    price DECIMAL(10, 2),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_by INT REFERENCES users(id),
    CONSTRAINT chk_copies CHECK (available_copies >= 0 AND available_copies <= total_copies),
    CONSTRAINT chk_publication_year CHECK (publication_year >= 1000 AND publication_year <= EXTRACT(YEAR FROM CURRENT_DATE) + 1)
);

-- Lendings Table
CREATE TABLE lendings (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE RESTRICT,
    status lending_status NOT NULL DEFAULT 'ACTIVE',
    issue_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    return_date TIMESTAMP WITH TIME ZONE,
    renewal_count INT DEFAULT 0,
    max_renewals INT DEFAULT 2,
    notes TEXT,
    issued_by INT REFERENCES users(id),
    returned_to INT REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_dates CHECK (due_date > issue_date),
    CONSTRAINT chk_return_date CHECK (return_date IS NULL OR return_date >= issue_date),
    CONSTRAINT chk_renewal_count CHECK (renewal_count >= 0 AND renewal_count <= max_renewals)
);

-- Reservations Table
CREATE TABLE reservations (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    book_id INT NOT NULL REFERENCES books(id) ON DELETE RESTRICT,
    status reservation_status NOT NULL DEFAULT 'PENDING',
    reservation_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expiry_date TIMESTAMP WITH TIME ZONE NOT NULL,
    fulfilled_date TIMESTAMP WITH TIME ZONE,
    cancelled_date TIMESTAMP WITH TIME ZONE,
    lending_id INT REFERENCES lendings(id),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_reservation_dates CHECK (expiry_date > reservation_date),
    CONSTRAINT chk_fulfilled_date CHECK (fulfilled_date IS NULL OR fulfilled_date >= reservation_date)
);

-- Fines Table
CREATE TABLE fines (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    lending_id INT REFERENCES lendings(id) ON DELETE SET NULL,
    amount DECIMAL(10, 2) NOT NULL,
    reason VARCHAR(255) NOT NULL,
    status fine_status NOT NULL DEFAULT 'PENDING',
    issued_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP WITH TIME ZONE NOT NULL,
    paid_date TIMESTAMP WITH TIME ZONE,
    payment_method VARCHAR(50),
    payment_reference VARCHAR(100),
    waived_by INT REFERENCES users(id),
    waived_date TIMESTAMP WITH TIME ZONE,
    waive_reason TEXT,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_amount CHECK (amount >= 0),
    CONSTRAINT chk_fine_dates CHECK (due_date >= issued_date),
    CONSTRAINT chk_paid_date CHECK (paid_date IS NULL OR paid_date >= issued_date)
);

-- Articles Table
CREATE TABLE articles (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    category article_category NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    author_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    featured_image_url TEXT,
    is_published BOOLEAN DEFAULT FALSE,
    published_at TIMESTAMP WITH TIME ZONE,
    views_count INT DEFAULT 0,
    tags TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_views CHECK (views_count >= 0)
);

-- Complaints Table
CREATE TABLE complaints (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    subject VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    category VARCHAR(100),
    priority VARCHAR(20) DEFAULT 'MEDIUM',
    status complaint_status NOT NULL DEFAULT 'OPEN',
    submitted_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    assigned_to INT REFERENCES users(id),
    resolved_date TIMESTAMP WITH TIME ZONE,
    resolution_notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_priority CHECK (priority IN ('LOW', 'MEDIUM', 'HIGH', 'URGENT')),
    CONSTRAINT chk_resolved_date CHECK (resolved_date IS NULL OR resolved_date >= submitted_date)
);

-- Create Indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_student_id ON users(student_id);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_users_status ON users(status);

CREATE INDEX idx_books_isbn ON books(isbn);
CREATE INDEX idx_books_title ON books(title);
CREATE INDEX idx_books_author ON books(author);
CREATE INDEX idx_books_category ON books(category);
CREATE INDEX idx_books_status ON books(status);

CREATE INDEX idx_lendings_user_id ON lendings(user_id);
CREATE INDEX idx_lendings_book_id ON lendings(book_id);
CREATE INDEX idx_lendings_status ON lendings(status);
CREATE INDEX idx_lendings_due_date ON lendings(due_date);
CREATE INDEX idx_lendings_issue_date ON lendings(issue_date);

CREATE INDEX idx_reservations_user_id ON reservations(user_id);
CREATE INDEX idx_reservations_book_id ON reservations(book_id);
CREATE INDEX idx_reservations_status ON reservations(status);

CREATE INDEX idx_fines_user_id ON fines(user_id);
CREATE INDEX idx_fines_status ON fines(status);
CREATE INDEX idx_fines_due_date ON fines(due_date);

CREATE INDEX idx_articles_slug ON articles(slug);
CREATE INDEX idx_articles_author_id ON articles(author_id);
CREATE INDEX idx_articles_category ON articles(category);
CREATE INDEX idx_articles_published_at ON articles(published_at);

CREATE INDEX idx_complaints_user_id ON complaints(user_id);
CREATE INDEX idx_complaints_status ON complaints(status);
CREATE INDEX idx_complaints_assigned_to ON complaints(assigned_to);

-- Create function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_books_updated_at BEFORE UPDATE ON books
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_lendings_updated_at BEFORE UPDATE ON lendings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_reservations_updated_at BEFORE UPDATE ON reservations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_fines_updated_at BEFORE UPDATE ON fines
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_articles_updated_at BEFORE UPDATE ON articles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_complaints_updated_at BEFORE UPDATE ON complaints
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Create function to auto-update book availability
CREATE OR REPLACE FUNCTION update_book_availability()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' AND NEW.status = 'ACTIVE' THEN
        UPDATE books 
        SET available_copies = available_copies - 1,
            status = CASE WHEN available_copies - 1 = 0 THEN 'BORROWED'::book_status ELSE status END
        WHERE id = NEW.book_id;
    ELSIF TG_OP = 'UPDATE' AND OLD.status = 'ACTIVE' AND NEW.status = 'RETURNED' THEN
        UPDATE books 
        SET available_copies = available_copies + 1,
            status = 'AVAILABLE'::book_status
        WHERE id = NEW.book_id;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_book_availability
AFTER INSERT OR UPDATE ON lendings
FOR EACH ROW EXECUTE FUNCTION update_book_availability();

-- Insert sample admin user (password: admin123 - bcrypt hashed)
INSERT INTO users (username, email, password_hash, first_name, last_name, role, employee_id, max_books_allowed)
VALUES ('admin', 'admin@library.com', '$2a$10$XQz5fJYCEU9ycPoQXZLYKuHCqLxZzXJXXB5rGx3qH5X5tYK5DQ5Ky', 'Admin', 'User', 'ADMIN', 'EMP001', 10);

-- Comments for documentation
COMMENT ON TABLE users IS 'Stores all system users including students, teachers, librarians, and admins';
COMMENT ON TABLE books IS 'Stores book inventory and metadata';
COMMENT ON TABLE lendings IS 'Tracks book borrowing and returns';
COMMENT ON TABLE reservations IS 'Manages book reservations when books are unavailable';
COMMENT ON TABLE fines IS 'Tracks fines for overdue books and other violations';
COMMENT ON TABLE articles IS 'Stores library news, events, and announcements';
COMMENT ON TABLE complaints IS 'Manages user complaints and feedback';