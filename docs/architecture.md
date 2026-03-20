## Architecture

This is a simple monlithic school library management system. It is built with Go (Golang) and Gin framework. This uses postgresql as the database. Here we have different roles for the users like Admin, Librarian, Teacher, Student. Each role has different permissions and capabilities. Admin has all the permissions and capabilities. Librarian has the permissions to manage the books, students, teachers, and other library related stuff. Teacher and student has the same level of permissions and capabilities in book lending and returning. In the article section, Teaches can disable and delete the article of the student but they can't submit the articles. Only Students can submit the article. To publish the article they need a review from the teacher. Teachers and student can rate and add comments in the article.

### Components

- **API**: The API is the main component of the application. It is responsible for handling the requests and responses from the client.
- **Database**: The database is responsible for storing the data of the application.
- **Authentication**: The authentication is responsible for handling the authentication of the users.
- **Authorization**: The authorization is responsible for handling the authorization of the users.
- **Logging**: The logging is responsible for logging the events of the application.
- **Monitoring**: The monitoring is responsible for monitoring the performance of the application.
- **Security**: The security is responsible for handling the security of the application.
- **Performance**: The performance is responsible for handling the performance of the application.
- **Scalability**: The scalability is responsible for handling the scalability of the application.
- **Availability**: The availability is responsible for handling the availability of the application.

### API Endpoints

- **Authentication**:
  - `POST /api/v1/auth/register` - Register a new user
  - `POST /api/v1/auth/login` - Login and get JWT token
  - `GET /api/v1/auth/profile` - Get current user profile
  - `POST /api/v1/auth/change-password` - Change password

- **Books Management**:
  - `GET /api/v1/books` - List all books (with pagination & search)
  - `GET /api/v1/books/:id` - Get book details
  - `POST /api/v1/books` - Create a new book (Admin/Librarian only)
  - `PUT /api/v1/books/:id` - Update a book (Admin/Librarian only)
  - `DELETE /api/v1/books/:id` - Delete a book (Admin/Librarian only)

- **Reviews**:
  - `GET /api/v1/books/:book_id/reviews` - List reviews for a book
  - `POST /api/v1/books/:book_id/reviews` - Add a review (Authenticated)
  - `PUT /api/v1/books/:book_id/reviews/:review_id` - Update a review (Owner/Admin/Librarian)
  - `DELETE /api/v1/books/:book_id/reviews/:review_id` - Delete a review (Owner/Admin/Librarian)

- **Comments**:
  - `GET /api/v1/books/:book_id/comments` - List comments for a book
  - `POST /api/v1/books/:book_id/comments` - Add a comment (Authenticated)
  - `PUT /api/v1/books/:book_id/comments/:comment_id` - Update a comment (Owner/Admin/Librarian)
  - `DELETE /api/v1/books/:book_id/comments/:comment_id` - Delete a comment (Owner/Admin/Librarian)

- **Articles**:
  - `GET /api/v1/articles` - List all articles (with pagination & search)
  - `GET /api/v1/articles/:id` - Get article details
  - `POST /api/v1/articles` - Create a new article (Student only)
  - `PUT /api/v1/articles/:id` - Update a article (Student only)
  - `DELETE /api/v1/articles/:id` - Delete a article (Student only)

- **Complaints**:
  - `GET /api/v1/complaints` - List all complaints (with pagination & search)
  - `GET /api/v1/complaints/:id` - Get complaint details
  - `POST /api/v1/complaints` - Create a new complaint (Authenticated)
  - `PUT /api/v1/complaints/:id` - Update a complaint (Owner/Admin/Librarian)
  - `DELETE /api/v1/complaints/:id` - Delete a complaint (Owner/Admin/Librarian)

- **Fines**:
  - `GET /api/v1/fines` - List all fines (with pagination & search)
  - `GET /api/v1/fines/:id` - Get fine details
  - `POST /api/v1/fines` - Create a new fine (Authenticated)

- **Lendings**:
  - `GET /api/v1/lendings` - List all lendings (with pagination & search)
  - `GET /api/v1/lendings/:id` - Get lending details
  - `POST /api/v1/lendings` - Create a new lending (Authenticated)
  - `PUT /api/v1/lendings/:id` - Update a lending (Owner/Admin/Librarian)
  - `DELETE /api/v1/lendings/:id` - Delete a lending (Owner/Admin/Librarian)

- **Reservations**:
  - `GET /api/v1/reservations` - List all reservations (with pagination & search)
  - `GET /api/v1/reservations/:id` - Get reservation details
  - `POST /api/v1/reservations` - Create a new reservation (Authenticated)
  - `PUT /api/v1/reservations/:id` - Update a reservation (Owner/Admin/Librarian)
  - `DELETE /api/v1/reservations/:id` - Delete a reservation (Owner/Admin/Librarian)

- **Users**:
  - `GET /api/v1/users` - List all users (with pagination & search)
  - `GET /api/v1/users/:id` - Get user details
  - `POST /api/v1/users` - Create a new user (Admin only)
  - `PUT /api/v1/users/:id` - Update a user (Admin only)
  - `DELETE /api/v1/users/:id` - Delete a user (Admin only)

- **Roles**:
  - `GET /api/v1/roles` - List all roles (with pagination & search)
  - `GET /api/v1/roles/:id` - Get role details
  - `POST /api/v1/roles` - Create a new role (Admin only)

- **Article Review**:
  - `GET /api/v1/articles/review` - List all article reviews (with pagination & search)
  - `GET /api/v1/articles/review/:id` - Get article review details
  - `POST /api/v1/articles/review` - Create a new article review (Teacher only)
  - `PUT /api/v1/articles/review/:id` - Update a article review (Teacher only)

- **Article Comments**:
  - `GET /api/v1/articles/:article_id/comments` - List comments for an article
  - `POST /api/v1/articles/:article_id/comments` - Add a comment (Authenticated)
  - `PUT /api/v1/articles/:article_id/comments/:comment_id` - Update a comment (Owner/Admin/Librarian)
  - `DELETE /api/v1/articles/:article_id/comments/:comment_id` - Delete a comment (Owner/Admin/Librarian)

- **Article Ratings**:
  - `GET /api/v1/articles/:article_id/ratings` - List ratings for an article
  - `POST /api/v1/articles/:article_id/ratings` - Add a rating (Authenticated)
  - `PUT /api/v1/articles/:article_id/ratings/:rating_id` - Update a rating (Owner/Admin/Librarian)
  - `DELETE /api/v1/articles/:article_id/ratings/:rating_id` - Delete a rating (Owner/Admin/Librarian)

