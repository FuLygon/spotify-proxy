package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spotify-proxy/internal/handlers"
)

type Routes interface {
	RegisterRoutes()
}

type routes struct {
	router            *gin.Engine
	proxyRouter       *gin.Engine
	authHandler       handlers.AuthHandler
	nowplayingHandler handlers.NowPlayingHandler
}

func NewRoutes(
	router *gin.Engine,
	proxyRouter *gin.Engine,
	authHandler handlers.AuthHandler,
	nowplayingHandler handlers.NowPlayingHandler,
) Routes {
	return &routes{
		router:            router,
		proxyRouter:       proxyRouter,
		authHandler:       authHandler,
		nowplayingHandler: nowplayingHandler,
	}
}

// RegisterRoutes configures API routes
func (r *routes) RegisterRoutes() {
	router := r.router
	proxyRouter := r.proxyRouter

	// Load HTML templates
	router.LoadHTMLGlob("static/*.html")

	// Serve index page
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{})
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Authentication endpoints
	authGroup := router.Group("/auth")
	authGroup.GET("/login", r.authHandler.HandleLogin)
	authGroup.GET("/callback", r.authHandler.HandleCallback)

	// Spotify NowPlayingHandler API
	proxyRouter.GET("/v1/me/player", r.nowplayingHandler.HandleCurrentTrack)
	proxyRouter.GET("/v1/me/player/queue", r.nowplayingHandler.HandleQueue)
}
