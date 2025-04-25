package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"spotify-proxy/config"
	"spotify-proxy/internal/services"
)

type Proxy interface {
	ReverseProxy(c *gin.Context)
}

type proxyHandler struct {
	config  *config.Config
	service services.AuthService
}

func NewProxyHandler(
	config *config.Config,
	service services.AuthService,
) Proxy {
	return &proxyHandler{
		config:  config,
		service: service,
	}
}

const spotifyAPIEndpoint = "https://api.spotify.com"

func (p *proxyHandler) ReverseProxy(c *gin.Context) {
	accessToken, err := p.service.GetAccessToken(c, p.config.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	target, err := url.Parse(spotifyAPIEndpoint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse target URL"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Director = func(req *http.Request) {
		req.Host = target.Host
		req.Header = c.Request.Header
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host

	}

	// Add custom error handler to handle CORS properly
	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		log.Printf("proxy error: %v\n", err)
		rw.WriteHeader(http.StatusBadGateway)
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
