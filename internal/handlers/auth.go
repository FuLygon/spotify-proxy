package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spotify-proxy/config"
	"spotify-proxy/internal/services"
)

type AuthHandler interface {
	HandleLogin(c *gin.Context)
	HandleCallback(c *gin.Context)
}

type authHandler struct {
	config  *config.Config
	service services.AuthService
}

func NewAuthHandler(config *config.Config, service services.AuthService) AuthHandler {
	return &authHandler{
		config:  config,
		service: service,
	}
}

func (h *authHandler) HandleLogin(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, h.service.Login(h.config.State))
}

func (h *authHandler) HandleCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")

	if state != h.config.State {
		c.JSON(http.StatusBadRequest, gin.H{"error": "state mismatch"})
		return
	}

	err := h.service.Callback(c, code, h.config.RefreshTokenOutput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Login successful"})
}
