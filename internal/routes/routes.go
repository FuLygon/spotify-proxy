package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spotify-proxy/config"
	"spotify-proxy/internal/handlers"
)

type Routes interface {
	RegisterRoutes()
}

type routes struct {
	router            *gin.Engine
	proxyRouter       *gin.Engine
	proxyRoutesConfig *config.ProxyRoutesConfig
	authHandler       handlers.AuthHandler
	proxyHandler      handlers.Proxy
}

func NewRoutes(
	router *gin.Engine,
	proxyRouter *gin.Engine,
	proxyRoutesConfig *config.ProxyRoutesConfig,
	authHandler handlers.AuthHandler,
	proxyHandler handlers.Proxy,
) Routes {
	return &routes{
		router:            router,
		proxyRouter:       proxyRouter,
		proxyRoutesConfig: proxyRoutesConfig,
		authHandler:       authHandler,
		proxyHandler:      proxyHandler,
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

	if r.proxyRoutesConfig != nil && len(r.proxyRoutesConfig.ProxyRoutes) > 0 {
		// Only proxy routes from the config file
		for _, route := range r.proxyRoutesConfig.ProxyRoutes {
			for _, method := range route.Methods {
				proxyRouter.Handle(method, route.Path, r.proxyHandler.ReverseProxy)
			}
		}
	} else {
		// Catch-all proxy
		proxyRouter.Any("/*path", r.proxyHandler.ReverseProxy)
	}
}
