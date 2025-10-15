package middleware

import (
	"library-management-system-go/pkg/response"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				response.Error(c, http.StatusInternalServerError, "Internal server error")
			}
		}()

		c.Next()

		// Check for errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error: %v", err.Error())
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
	}
}
