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

	reviews := books.Group("/:id/reviews")
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

	comments := books.Group("/:id/comments")
	{
		comments.GET("", commentHandler.ListComments)

		// Authenticated users can add comments
		comments.POST("", middleware.AuthMiddleware(authService), commentHandler.CreateComment)

		// Owner, Admin, Librarian can update/delete
		comments.PUT("/:comment_id", middleware.AuthMiddleware(authService), commentHandler.UpdateComment)
		comments.DELETE("/:comment_id", middleware.AuthMiddleware(authService), commentHandler.DeleteComment)
	}

	// Articles and article sub-resources
	articleRepo := repository.NewArticleRepository(db)
	articleReviewRepo := repository.NewArticleReviewRepository(db)
	articleCommentRepo := repository.NewArticleCommentRepository(db)
	articleRatingRepo := repository.NewArticleRatingRepository(db)

	articleService := service.NewArticleService(articleRepo, articleReviewRepo, articleCommentRepo, articleRatingRepo)
	articleHandler := handler.NewArticleHandler(articleService)

	articles := protected.Group("/articles")
	{
		articles.GET("", articleHandler.ListArticles)
		articles.GET("/:id", articleHandler.GetArticle)

		articles.POST("", middleware.RoleMiddleware("STUDENT"), articleHandler.CreateArticle)
		articles.PUT("/:id", middleware.RoleMiddleware("STUDENT", "ADMIN"), articleHandler.UpdateArticle)
		articles.DELETE("/:id", middleware.RoleMiddleware("STUDENT", "TEACHER", "ADMIN"), articleHandler.DeleteArticle)
	}

	articleReviews := protected.Group("/articles/review")
	{
		articleReviews.Use(middleware.RoleMiddleware("TEACHER", "ADMIN"))
		{
			articleReviews.GET("", articleHandler.ListArticleReviews)
			articleReviews.GET("/:id", articleHandler.GetArticleReview)
			articleReviews.POST("", articleHandler.CreateArticleReview)
			articleReviews.PUT("/:id", articleHandler.UpdateArticleReview)
		}
	}

	articleComments := protected.Group("/articles/:id/comments")
	{
		articleComments.GET("", articleHandler.ListArticleComments)
		articleComments.POST("", articleHandler.CreateArticleComment)
		articleComments.PUT("/:comment_id", articleHandler.UpdateArticleComment)
		articleComments.DELETE("/:comment_id", articleHandler.DeleteArticleComment)
	}

	articleRatings := protected.Group("/articles/:id/ratings")
	{
		articleRatings.GET("", articleHandler.ListArticleRatings)
		articleRatings.POST("", articleHandler.CreateArticleRating)
		articleRatings.PUT("/:rating_id", articleHandler.UpdateArticleRating)
		articleRatings.DELETE("/:rating_id", articleHandler.DeleteArticleRating)
	}

	// Lendings
	lendingRepo := repository.NewLendingRepository(db)
	lendingService := service.NewLendingService(lendingRepo)
	lendingHandler := handler.NewLendingHandler(lendingService)

	lendings := protected.Group("/lendings")
	{
		lendings.GET("", lendingHandler.ListLendings)
		lendings.GET("/:id", lendingHandler.GetLending)
		lendings.POST("", lendingHandler.CreateLending)
		lendings.PUT("/:id", lendingHandler.UpdateLending)
		lendings.DELETE("/:id", lendingHandler.DeleteLending)
	}

	// Reservations
	reservationRepo := repository.NewReservationRepository(db)
	reservationService := service.NewReservationService(reservationRepo)
	reservationHandler := handler.NewReservationHandler(reservationService)

	reservations := protected.Group("/reservations")
	{
		reservations.GET("", reservationHandler.ListReservations)
		reservations.GET("/:id", reservationHandler.GetReservation)
		reservations.POST("", reservationHandler.CreateReservation)
		reservations.PUT("/:id", reservationHandler.UpdateReservation)
		reservations.DELETE("/:id", reservationHandler.DeleteReservation)
	}

	// Fines
	fineRepo := repository.NewFineRepository(db)
	fineService := service.NewFineService(fineRepo, lendingRepo)
	fineHandler := handler.NewFineHandler(fineService)

	fines := protected.Group("/fines")
	{
		fines.GET("", fineHandler.ListFines)
		fines.GET("/:id", fineHandler.GetFine)
		fines.POST("", fineHandler.CreateFine)
	}

	// Complaints
	complaintRepo := repository.NewComplaintRepository(db)
	complaintService := service.NewComplaintService(complaintRepo)
	complaintHandler := handler.NewComplaintHandler(complaintService)

	complaints := protected.Group("/complaints")
	{
		complaints.GET("", complaintHandler.ListComplaints)
		complaints.GET("/:id", complaintHandler.GetComplaint)
		complaints.POST("", complaintHandler.CreateComplaint)
		complaints.PUT("/:id", complaintHandler.UpdateComplaint)
		complaints.DELETE("/:id", complaintHandler.DeleteComplaint)
	}

	// Users (admin-only)
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	users := protected.Group("/users")
	{
		users.Use(middleware.RoleMiddleware("ADMIN"))
		{
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// Roles (fixed system enum)
	roleService := service.NewRoleService()
	roleHandler := handler.NewRoleHandler(roleService)

	roles := protected.Group("/roles")
	{
		roles.GET("", roleHandler.ListRoles)
		roles.GET("/:id", roleHandler.GetRole)
		roles.POST("", middleware.RoleMiddleware("ADMIN"), roleHandler.CreateRole)
	}
}
