package routes

import (
	"github.com/gin-gonic/gin"
	"spotify-proxy/internal/handlers"
)

type Routes interface {
	RegisterRoutes()
}

type routes struct {
	router       *gin.Engine
	proxyRouter  *gin.Engine
	authHandler  handlers.AuthHandler
	proxyHandler handlers.Proxy
}

func NewRoutes(
	router *gin.Engine,
	proxyRouter *gin.Engine,
	authHandler handlers.AuthHandler,
	proxyHandler handlers.Proxy,
) Routes {
	return &routes{
		router:       router,
		proxyRouter:  proxyRouter,
		authHandler:  authHandler,
		proxyHandler: proxyHandler,
	}
}

// RegisterRoutes configures API routes
func (r *routes) RegisterRoutes() {
	router := r.router
	proxyRouter := r.proxyRouter

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Authentication endpoints
	authGroup := router.Group("/auth")
	authGroup.GET("/login", r.authHandler.HandleLogin)
	authGroup.GET("/callback", r.authHandler.HandleCallback)

	// Reverse proxy endpoints
	proxyRouter.Any("/*path", r.proxyHandler.ReverseProxy)
}
