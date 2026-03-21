# Postman Collection Setup Guide

## Overview
This directory contains a complete Postman collection for testing all endpoints of the School Library Management API. All default test users are included with their credentials.

## Files
- `postman_collection.json` - Complete API collection with all endpoints (34 paths, 64 methods)
- `postman_environment.json` - Environment variables including base URL and default credentials

## Import Instructions

### 1. Import Collection
1. Open Postman
2. Click **File** → **Import** (or use `Ctrl+O`)
3. Select **File** tab
4. Browse and select `postman_collection.json`
5. Click **Import**

### 2. Import Environment
1. Click the **gear icon** (Settings) in top-right
2. Select **Manage Environments**
3. Click **Import**
4. Browse and select `postman_environment.json`
5. Click **Import**
6. Close the modal and select the imported environment from the dropdown

## Default Test Credentials

All of these users are created automatically when you run the seeder:

| Username | Password | Role | Email |
|----------|----------|------|-------|
| admin | admin123 | Admin | admin@library.com |
| admin_2 | admin29 | Admin | admin_2@library.com |
| librarian | librarian123 | Librarian | librarian@library.com |
| teacher | teacher123 | Teacher | teacher@library.com |
| student | student123 | Student | student@library.com |

## Quick Start

### Step 1: Login
1. Click **Authentication** folder
2. Click **Login** request
3. Edit the body to use credentials from above (pre-filled with admin)
4. Click **Send**
5. The token will automatically be saved to the `token` variable

### Step 2: Test API Endpoints
1. Browse the collection folders by resource type
2. Click any request and modify parameters/body as needed
3. Click **Send** to execute
4. View response in the Response panel

## Collection Structure

```
├── Authentication
│   ├── Register
│   ├── Login (auto-saves token)
│   ├── Get Profile
│   └── Change Password
├── Role Access
│   ├── Admin Hello
│   ├── Librarian Hello
│   ├── Student Hello
│   └── Teacher Hello
├── Books
│   ├── List Books
│   ├── Get Book
│   ├── Create Book
│   ├── Update Book
│   └── Delete Book
├── Book Reviews
│   ├── List Reviews
│   ├── Create Review
│   ├── Update Review
│   └── Delete Review
├── Book Comments
│   ├── List Comments
│   ├── Create Comment
│   ├── Update Comment
│   └── Delete Comment
├── Articles
│   ├── List Articles
│   ├── Get Article
│   ├── Create Article
│   ├── Update Article
│   └── Delete Article
├── Article Reviews
│   ├── List Article Reviews
│   ├── Get Article Review
│   ├── Create Article Review
│   └── Update Article Review
├── Article Comments
│   ├── List Article Comments
│   ├── Create Article Comment
│   ├── Update Article Comment
│   └── Delete Article Comment
├── Article Ratings
│   ├── List Article Ratings
│   ├── Create Article Rating
│   ├── Update Article Rating
│   └── Delete Article Rating
├── Lendings
│   ├── List Lendings
│   ├── Get Lending
│   ├── Create Lending
│   ├── Update Lending
│   └── Delete Lending
├── Reservations
│   ├── List Reservations
│   ├── Get Reservation
│   ├── Create Reservation
│   ├── Update Reservation
│   └── Delete Reservation
├── Fines
│   ├── List Fines
│   ├── Get Fine
│   └── Create Fine
├── Complaints
│   ├── List Complaints
│   ├── Get Complaint
│   ├── Create Complaint
│   ├── Update Complaint
│   └── Delete Complaint
├── Users (Admin only)
│   ├── List Users
│   ├── Get User
│   ├── Create User
│   ├── Update User
│   └── Delete User
└── Roles
    ├── List Roles
    ├── Get Role
    └── Create Role
```

## Important Notes

### Authentication
- The **Login** request includes a test script that automatically extracts the JWT token and saves it to the `{{token}}` variable
- All other requests use Bearer token authentication via the collection's auth settings
- Token will expire; re-login if you get 401 Unauthorized

### Example Request Bodies
- All POST/PUT requests include example payloads
- Modify the request bodies before sending to match your test data
- Check the example payloads for field names and data types

### Pagination
- List endpoints support `page`, `limit`, and `search` parameters
- Default: page=1, limit=10
- Modify these in the URL query parameters as needed

### Role-Based Access
- Some endpoints are restricted by role (e.g., Users, Admin/Librarian routes)
- Test role-specific access by logging in with different users
- Use the Role Access folder to test role-based endpoint access

## Testing Workflow

### Full Flow Test
1. **Login** as a specific role (e.g., admin, student, librarian)
2. **Create** a resource (e.g., Create Article)
3. **List** resources to verify creation
4. **Get** the specific resource
5. **Update** the resource
6. **Delete** the resource

### Example: Testing Book Flow
1. Login as librarian (has book management permissions)
2. Create a book
3. List books to find your new book
4. Get the specific book details
5. Create a review on the book
6. List reviews
7. Update the review
8. Delete the review
9. Update the book
10. Delete the book

## Variables

The environment file includes pre-configured variables:

- `{{base_url}}` - API base URL (http://localhost:8080/api/v1)
- `{{token}}` - JWT authentication token (auto-set after login)
- `{{admin_username}}`, `{{admin_password}}` - Admin credentials
- `{{librarian_username}}`, `{{librarian_password}}` - Librarian credentials
- `{{teacher_username}}`, `{{teacher_password}}` - Teacher credentials
- `{{student_username}}`, `{{student_password}}` - Student credentials

## Troubleshooting

### 401 Unauthorized
- Token has expired
- Solution: Run the **Login** request again to refresh the token

### 403 Forbidden
- Your role doesn't have permission for this endpoint
- Solution: Login with a user that has the required role

### 404 Not Found
- Resource ID doesn't exist
- Solution: Use valid resource IDs from List endpoints

### CORS Errors
- Ensure the API server is running on http://localhost:8080
- Check that CORS middleware is properly configured

### Request Variables Not Working
- Ensure the environment is selected (dropdown in top-right)
- Check that variable names use `{{variable}}` syntax

## API Documentation

For detailed endpoint documentation, visit:
- Swagger UI: http://localhost:8080/swagger/index.html
- Swagger JSON: http://localhost:8080/swagger/doc.json

## Notes

- All timestamps are in ISO 8601 format
- Dates should be in YYYY-MM-DD format
- Passwords are hashed using bcrypt with salt
- User emails and usernames must be unique
- Maximum of 5 books for students, 10 for teachers, 15 for librarians, 20 for admins
