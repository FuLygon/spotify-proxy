package main

import (
	"github.com/gin-gonic/gin"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
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

	// Setup spotify authenticator
	sa := spotifyAuth.New(
		spotifyAuth.WithRedirectURL(conf.RedirectUri),
		spotifyAuth.WithClientID(conf.ClientId),
		spotifyAuth.WithClientSecret(conf.ClientSecret),
		spotifyAuth.WithScopes(conf.Scope...),
	)

	// Setup services
	authService := services.NewAuthService(sa, cacheInstance)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(conf, authService)

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
