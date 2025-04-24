package services

import (
	"spotify-proxy/internal/cache"
)

type AuthService interface {
	Callback()
}

type authService struct {
	cache cache.Cache
}

func NewAuthService(cache cache.Cache) AuthService {
	return &authService{
		cache: cache,
	}
}

func (s *authService) Callback() {
	return
}
