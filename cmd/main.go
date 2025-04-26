package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"spotify-proxy/config"
	"spotify-proxy/internal/cache"
	"spotify-proxy/internal/handlers"
	"spotify-proxy/internal/routes"
	"spotify-proxy/internal/services"
	"time"
)

var g errgroup.Group

func main() {
	// Load env
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load configuration: ", err)
	}

	// Set Gin mode
	gin.SetMode(conf.GinMode)
	accessRouter := gin.Default()
	nowplayingRouter := gin.Default()

	// Now Playing CORS
	corsConfig := cors.DefaultConfig()
	if len(conf.CorsOrigins) > 0 {
		corsConfig.AllowOrigins = conf.CorsOrigins
	} else {
		corsConfig.AllowAllOrigins = true
	}
	nowplayingRouter.Use(cors.New(corsConfig))

	// Set trusted proxies
	if len(conf.TrustedProxies) > 0 {
		err = accessRouter.SetTrustedProxies(conf.TrustedProxies)
		if err != nil {
			log.Printf("error setting trusted proxies: %v", err)
		}

		err = nowplayingRouter.SetTrustedProxies(conf.TrustedProxies)
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
	nowplayingService := services.NewNowPlayingService(sa, cacheInstance)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(conf, authService)
	nowPlayingHandler := handlers.NewNowPlayingHandler(conf, authService, nowplayingService)

	// Setup routes
	r := routes.NewRoutes(
		accessRouter,
		nowplayingRouter,
		authHandler,
		nowPlayingHandler,
	)
	if err != nil {
		log.Fatal("failed to setup routes: ", err)
	}

	// Register routes
	r.RegisterRoutes()

	accessServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", conf.AccessPort),
		Handler:      accessRouter.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	nowplayingServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", conf.NowPlayingPort),
		Handler:      nowplayingRouter.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		log.Printf("Access server is listening on %s", accessServer.Addr)
		return accessServer.ListenAndServe()
	})

	g.Go(func() error {
		log.Printf("Now Playing server is listening on %s", nowplayingServer.Addr)
		return nowplayingServer.ListenAndServe()
	})

	if err = g.Wait(); err != nil {
		log.Fatal(err)
	}
}
