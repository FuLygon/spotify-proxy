package routes

import (
	"github.com/gin-gonic/gin"
	"spotify-proxy/internal/handlers"
)

type Routes interface {
	RegisterRoutes()
}

type routes struct {
	router      *gin.Engine
	authHandler handlers.AuthHandler
}

func NewRoutes(
	router *gin.Engine,
	authHandler handlers.AuthHandler,
) Routes {
	return &routes{
		router:      router,
		authHandler: authHandler,
	}
}

// RegisterRoutes configures API routes
func (r *routes) RegisterRoutes() {
	router := r.router

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	authGroup := r.router.Group("/auth")
	authGroup.GET("/login", r.authHandler.HandleLogin)
	authGroup.GET("/callback", r.authHandler.HandleCallback)
}
