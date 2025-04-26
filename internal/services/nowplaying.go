package services

import (
	"encoding/json"
	"fmt"
	spotifyAuth "github.com/zmb3/spotify/v2/auth"
	"net/http"
	"spotify-proxy/internal/cache"
	"spotify-proxy/internal/models"
	"time"
)

type NowPlayingService interface {
	GetPlayerCurrentTrack(currentTrack, accessToken string) (map[string]interface{}, error)
	GetRecentlyPlayedTracks(endpoint, accessToken string) (map[string]interface{}, error)
	GetPlayerQueue(endpoint, accessToken string, cacheInterval int) (*models.PlayerQueue, error)
}

type nowplayingService struct {
	sa     *spotifyAuth.Authenticator
	cache  cache.Cache
	client *http.Client
}

const currentQueueCacheKey = "current_queue"

var currentQueueCacheLastUpdate time.Time

func NewNowPlayingService(sa *spotifyAuth.Authenticator, cache cache.Cache) NowPlayingService {
	return &nowplayingService{
		sa:    sa,
		cache: cache,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *nowplayingService) GetPlayerCurrentTrack(currentTrack, accessToken string) (map[string]interface{}, error) {
	statsReq, err := http.NewRequest(http.MethodGet, currentTrack, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare stats request: %w", err)
	}

	statsReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Make stats request
	resp, err := s.client.Do(statsReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, nil
	}

	// Parse stats response
	var response map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse stats response: %w", err)
	}

	return response, nil
}

func (s *nowplayingService) GetRecentlyPlayedTracks(endpoint, accessToken string) (map[string]interface{}, error) {
	statsReq, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare stats request: %w", err)
	}

	statsReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Make stats request
	resp, err := s.client.Do(statsReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	defer resp.Body.Close()

	// Parse stats response
	var recentTrack models.PlayerRecentTrack
	if err = json.NewDecoder(resp.Body).Decode(&recentTrack); err != nil {
		return nil, fmt.Errorf("failed to parse stats response: %w", err)
	}

	if len(recentTrack.Items) > 0 {
		response := make(map[string]interface{})
		response["item"] = recentTrack.Items[0].Track
		response["is_playing"] = false
		return response, nil
	} else {
		return nil, nil
	}
}

func (s *nowplayingService) GetPlayerQueue(endpoint, accessToken string, cacheInterval int) (*models.PlayerQueue, error) {
	statsReq, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare stats request: %w", err)
	}

	statsReq.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	// Make stats request
	resp, err := s.client.Do(statsReq)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch stats: %w", err)
	}
	defer resp.Body.Close()

	// Parse stats response
	var response models.PlayerQueue
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse stats response: %w", err)
	}

	if len(response.Queue) == 0 {
		currentQueueCache, found := s.cache.Get(currentQueueCacheKey)
		if found {
			// return cached response if player is not playing
			response = currentQueueCache.(models.PlayerQueue)
			return &response, nil
		} else {
			// return empty response
			dummyResponse := make([]map[string]interface{}, 0)
			return &models.PlayerQueue{Queue: dummyResponse}, nil
		}
	}

	// Cache the response after specific interval
	if time.Since(currentQueueCacheLastUpdate) > time.Duration(cacheInterval)*time.Second {
		s.cache.Set(currentQueueCacheKey, response, cache.NoExpiration)
		currentQueueCacheLastUpdate = time.Now()
	}

	return &response, nil
}
