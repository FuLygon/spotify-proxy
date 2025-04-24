package services

import (
	"context"
	"fmt"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
	"spotify-proxy/internal/cache"
	"time"
)

type AuthService interface {
	Login(state string) string
	Callback(ctx context.Context, code string, refreshTokenLog bool) error
}

type authService struct {
	sa    *spotifyAuth.Authenticator
	cache cache.Cache
}

const (
	accessTokenCacheKey  = "access_token"
	refreshTokenCacheKey = "refresh_token"
)

func NewAuthService(sa *spotifyAuth.Authenticator, cache cache.Cache) AuthService {
	return &authService{
		sa:    sa,
		cache: cache,
	}
}

func (s *authService) Login(state string) string {
	return s.sa.AuthURL(state)
}

func (s *authService) Callback(ctx context.Context, code string, refreshTokenLog bool) error {
	token, err := s.sa.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("token exchange error: %w", err)
	}

	// Refresh token cache will have no expiration
	s.cache.Set(refreshTokenCacheKey, token.RefreshToken, cache.NoExpiration)

	// Access token cache will expire 1 minute before the actual expiration time
	s.cache.Set(accessTokenCacheKey, token.AccessToken, token.Expiry.Sub(time.Now())-time.Minute)

	// Optionally log the refresh token
	if refreshTokenLog {
		fmt.Printf("\nRefresh token: %s\nYou can set this value into SPOTIFY_REFRESH_TOKEN so you don't have to login next time on startup\n\n", token.RefreshToken)
	}
	return nil
}
