package main

import (
	"fmt"
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

	// Load proxy route config
	proxyRoutesConf, err := config.LoadProxyRoutesConfig("./routes.yaml")
	if err != nil {
		log.Fatal("failed to load proxy route configuration: ", err)
	}

	// Set Gin mode
	gin.SetMode(conf.GinMode)
	accessRouter := gin.Default()
	proxyRouter := gin.Default()

	// Set trusted proxies
	if len(conf.TrustedProxies) > 0 {
		err = accessRouter.SetTrustedProxies(conf.TrustedProxies)
		if err != nil {
			log.Printf("error setting trusted proxies: %v", err)
		}

		err = proxyRouter.SetTrustedProxies(conf.TrustedProxies)
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
	proxyHandler := handlers.NewProxyHandler(conf, authService)

	// Setup routes
	r := routes.NewRoutes(
		accessRouter,
		proxyRouter,
		proxyRoutesConf,
		authHandler,
		proxyHandler,
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

	proxyServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", conf.ProxyPort),
		Handler:      proxyRouter.Handler(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		log.Printf("Access server is listening on %s", accessServer.Addr)
		return accessServer.ListenAndServe()
	})

	g.Go(func() error {
		log.Printf("Proxy server is listening on %s", proxyServer.Addr)
		return proxyServer.ListenAndServe()
	})

	if err = g.Wait(); err != nil {
		log.Fatal(err)
	}
}
