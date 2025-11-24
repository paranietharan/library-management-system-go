package api

import (
	"library-management-system-go/internal/config"
	"library-management-system-go/internal/handler"
	"library-management-system-go/internal/middleware"
	"library-management-system-go/internal/repository"
	"library-management-system-go/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, authService service.AuthService, cfg *config.Config, db *gorm.DB) {
	authHandler := handler.NewAuthHandler(authService)

	auth := router.Group("/api/v1/auth")
	auth.Use(middleware.CORSMiddleware(cfg))
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		protected.GET("/auth/profile", authHandler.GetProfile)
		protected.POST("/auth/change-password", authHandler.ChangePassword)

		admin := protected.Group("/admin")
		admin.Use(middleware.RoleMiddleware("ADMIN"))
		{
			admin.GET("/hello", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Hello Admin! You have successfully passed authentication and role verification.",
				})
			})
		}

		librarian := protected.Group("/librarian")
		librarian.Use(middleware.RoleMiddleware("ADMIN", "LIBRARIAN"))
		{
			librarian.GET("/hello", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Hello Librarian (or Admin)! Your access is verified.",
				})
			})
		}

		student := protected.Group("/student")
		student.Use(middleware.RoleMiddleware("STUDENT"))
		{
			student.GET("/hello", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Hello Student! Your access is verified.",
				})
			})
		}

		teacher := protected.Group("/teacher")
		teacher.Use(middleware.RoleMiddleware("TEACHER"))
		{
			teacher.GET("/hello", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "Hello Teacher! Your access is verified.",
				})
			})
		}
	}

	// Books management
	bookRepo := repository.NewBookRepository(db)
	bookService := service.NewBookService(bookRepo)
	bookHandler := handler.NewBookHandler(bookService)

	books := protected.Group("/books")
	{
		books.GET("", bookHandler.ListBooks)
		books.GET("/:id", bookHandler.GetBook)

		// Admin and Librarian only
		manage := books.Group("")
		manage.Use(middleware.RoleMiddleware("ADMIN", "LIBRARIAN"))
		{
			manage.POST("", bookHandler.CreateBook)
			manage.PUT("/:id", bookHandler.UpdateBook)
			manage.DELETE("/:id", bookHandler.DeleteBook)
		}
	}

	// Endpoints related to books rating and book feed
	reviewRepo := repository.NewReviewRepository(db)
	reviewService := service.NewReviewService(reviewRepo)
	reviewHandler := handler.NewReviewHandler(reviewService)

	reviews := books.Group("/:book_id/reviews")
	{
		reviews.GET("", reviewHandler.ListReviews)

		// Authenticated users can add reviews
		reviews.POST("", middleware.AuthMiddleware(authService), reviewHandler.CreateReview)

		// Owner, Admin, Librarian can update/delete
		reviews.PUT("/:review_id", middleware.AuthMiddleware(authService), reviewHandler.UpdateReview)
		reviews.DELETE("/:review_id", middleware.AuthMiddleware(authService), reviewHandler.DeleteReview)
	}

	commentRepo := repository.NewCommentRepository(db)
	commentService := service.NewCommentService(commentRepo)
	commentHandler := handler.NewCommentHandler(commentService)

	comments := books.Group("/:book_id/comments")
	{
		comments.GET("", commentHandler.ListComments)

		// Authenticated users can add comments
		comments.POST("", middleware.AuthMiddleware(authService), commentHandler.CreateComment)

		// Owner, Admin, Librarian can update/delete
		comments.PUT("/:comment_id", middleware.AuthMiddleware(authService), commentHandler.UpdateComment)
		comments.DELETE("/:comment_id", middleware.AuthMiddleware(authService), commentHandler.DeleteComment)
	}
}
