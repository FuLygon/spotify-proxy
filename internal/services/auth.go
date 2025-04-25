package services

import (
	"context"
	"fmt"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
	"spotify-proxy/internal/cache"
	"time"
)

type AuthService interface {
	Login(state string) string
	Callback(ctx context.Context, code string, refreshTokenLog bool) error
	GetAccessToken(ctx context.Context, refreshTokenEnv string) (string, error)
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

func (s *authService) GetAccessToken(ctx context.Context, refreshTokenEnv string) (string, error) {
	// Get access token from cache
	accessTokenCached, found := s.cache.Get(accessTokenCacheKey)
	if found {
		return accessTokenCached.(string), nil
	}

	// Access token either not found or expired, attempt to acquire a new one from refresh token
	var (
		token *oauth2.Token
		err   error
	)

	if refreshTokenEnv != "" {
		// Use refresh token from env if available
		token, err = s.refreshToken(ctx, &oauth2.Token{RefreshToken: refreshTokenEnv})
		if err != nil {
			return "", fmt.Errorf("failed to refresh token from env: %w", err)
		}
	} else {
		// Attempt to get refresh token from cache
		refreshTokenCached, found := s.cache.Get(refreshTokenCacheKey)
		if found {
			token, err = s.refreshToken(ctx, &oauth2.Token{RefreshToken: refreshTokenCached.(string)})
			if err != nil {
				return "", fmt.Errorf("failed to refresh token from env: %w", err)
			}
		} else {
			return "", fmt.Errorf("no refresh token found in cache or env, try logging in")
		}
	}

	// Access token cache will expire 1 minute before the actual expiration time
	s.cache.Set(accessTokenCacheKey, token.AccessToken, token.Expiry.Sub(time.Now())-time.Minute)
	return token.AccessToken, nil
}

func (s *authService) refreshToken(ctx context.Context, refreshToken *oauth2.Token) (*oauth2.Token, error) {
	token, err := s.sa.RefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return token, nil
}
