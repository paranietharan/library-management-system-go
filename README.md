# School Library Management System

A robust backend system for managing a school library, built with Go (Golang) and Gin framework.

## Features

- **Authentication & Authorization**: JWT-based auth with role-based access control (Admin, Librarian, Teacher, Student).
- **User Management**: Manage users with different roles.
- **Book Management**: Complete CRUD operations for books with inventory tracking.
- **PostgreSQL Database**: Reliable data storage using GORM.

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authentication**: JWT (JSON Web Tokens)
- **Configuration**: Godotenv

## Setup & Installation

1.  **Clone the repository**
    ```bash
    git clone <repository-url>
    cd library-management-system-go
    ```

2.  **Environment Setup**
    Copy `.env.example` to `.env` and update the values:
    ```bash
    cp .env.example .env
    ```

3.  **Run the Application**
    ```bash
    go run cmd/run/main.go
    ```

## API Documentation (Swagger)

This project now exposes Swagger UI for interactive API docs.

- Open: `http://localhost:8080/swagger/index.html`
- JSON spec: `http://localhost:8080/swagger/doc.json`

To regenerate docs after API comment changes:
```bash
go install github.com/swaggo/swag/cmd/swag@latest
"$(go env GOPATH)/bin/swag" init -g cmd/run/main.go -o docs/swagger
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token
- `GET /api/v1/auth/profile` - Get current user profile
- `POST /api/v1/auth/change-password` - Change password

### Books Management
- `GET /api/v1/books` - List all books (with pagination & search)
- `GET /api/v1/books/:id` - Get book details
- `POST /api/v1/books` - Create a new book (Admin/Librarian only)
- `PUT /api/v1/books/:id` - Update a book (Admin/Librarian only)
- `DELETE /api/v1/books/:id` - Delete a book (Admin/Librarian only)

### Reviews
- `GET /api/v1/books/:id/reviews` - List reviews for a book
- `POST /api/v1/books/:id/reviews` - Add a review (Authenticated)
- `PUT /api/v1/books/:id/reviews/:review_id` - Update a review (Owner/Admin/Librarian)
- `DELETE /api/v1/books/:id/reviews/:review_id` - Delete a review (Owner/Admin/Librarian)

### Comments
- `GET /api/v1/books/:id/comments` - List comments for a book
- `POST /api/v1/books/:id/comments` - Add a comment (Authenticated)
- `PUT /api/v1/books/:id/comments/:comment_id` - Update a comment (Owner/Admin/Librarian)
- `DELETE /api/v1/books/:id/comments/:comment_id` - Delete a comment (Owner/Admin/Librarian)

## Testing

Run unit tests:
```bash
go test ./...
```
