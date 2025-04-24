package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"spotify-proxy/config"
	"spotify-proxy/internal/cache"
	"spotify-proxy/internal/handlers"
	"spotify-proxy/internal/routes"
	"spotify-proxy/internal/services"
)

func main() {
	// Load env
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load configuration: ", err)
	}

	// Set Gin mode
	gin.SetMode(conf.GinMode)
	router := gin.Default()

	// Set trusted proxies
	if len(conf.TrustedProxies) > 0 {
		err = router.SetTrustedProxies(conf.TrustedProxies)
		if err != nil {
			log.Printf("error setting trusted proxies: %v", err)
		}
	}

	// Initialize cache
	cacheInstance := cache.NewCache()

	// Setup services
	authService := services.NewAuthService(cacheInstance)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(authService)

	// Setup routes
	r := routes.NewRoutes(
		router,
		authHandler,
	)

	// Register routes
	r.RegisterRoutes()

	// Start server
	if err = router.Run(":" + conf.Port); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
