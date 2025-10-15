package api

import (
	"library-management-system-go/internal/config"
	"library-management-system-go/internal/handler"
	"library-management-system-go/internal/middleware"
	"library-management-system-go/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, authService service.AuthService, cfg *config.Config) {
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
}
