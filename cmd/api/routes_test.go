package api

import (
	"library-management-system-go/internal/config"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSetupRoutes_NoGinRouteConflict(t *testing.T) {
	router := gin.New()

	cfg := &config.Config{
		Security: config.SecurityConfig{
			AllowedOrigins: "*",
		},
		Server: config.ServerConfig{
			Env: "development",
		},
		JWT: config.JWTConfig{
			Secret:      "test-secret",
			ExpiryHours: 1,
		},
	}

	// This test intentionally does not make any HTTP requests, it only ensures route registration
	// does not panic due to Gin wildcard conflicts.
	assert.NotPanics(t, func() {
		SetupRoutes(router, nil, cfg, &gorm.DB{})
	})
}

