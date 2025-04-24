package handlers

import (
	"github.com/gin-gonic/gin"
	"spotify-proxy/internal/services"
)

type AuthHandler interface {
	HandleLogin(c *gin.Context)
	HandleCallback(c *gin.Context)
}

type authHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) AuthHandler {
	return &authHandler{
		service: service,
	}
}

func (h *authHandler) HandleLogin(c *gin.Context) {
	c.JSON(200, nil)
}

func (h *authHandler) HandleCallback(c *gin.Context) {
	c.JSON(200, nil)
}
