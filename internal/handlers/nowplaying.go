package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"spotify-proxy/config"
	"spotify-proxy/internal/services"
)

type NowPlayingHandler interface {
	HandleCurrentTrack(c *gin.Context)
	HandleQueue(c *gin.Context)
}

type nowplayingHandler struct {
	config        *config.Config
	authService   services.AuthService
	playerService services.NowPlayingService
}

func NewNowPlayingHandler(
	config *config.Config,
	authService services.AuthService,
	playerService services.NowPlayingService,
) NowPlayingHandler {
	return &nowplayingHandler{
		config:        config,
		authService:   authService,
		playerService: playerService,
	}
}

const spotifyAPIEndpoint = "https://api.spotify.com"

func (p *nowplayingHandler) HandleCurrentTrack(c *gin.Context) {
	accessToken, err := p.authService.GetAccessToken(c, p.config.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Attempt to get currently playing track
	currentTrackReqUrl, err := url.Parse(fmt.Sprintf("%s/v1/me/player", spotifyAPIEndpoint))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to parse URL"})
		return
	}
	currentTrackReqUrl.RawQuery = c.Request.URL.Query().Encode()

	resp, err := p.playerService.GetPlayerCurrentTrack(currentTrackReqUrl.String(), accessToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Recently played track will be used if the current track is nil
	recentTrackReqUrl, err := url.Parse(fmt.Sprintf("%s/v1/me/player/recently-played?limit=1", spotifyAPIEndpoint))
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to parse URL"})
		return
	}

	if resp == nil {
		resp, err = p.playerService.GetRecentlyPlayedTracks(recentTrackReqUrl.String(), accessToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else if resp == nil {
			c.JSON(http.StatusNoContent, nil)
			return
		}
	}

	c.JSON(http.StatusOK, resp)

}

func (p *nowplayingHandler) HandleQueue(c *gin.Context) {
	accessToken, err := p.authService.GetAccessToken(c, p.config.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	reqUrl := fmt.Sprintf("%s/v1/me/player/queue", spotifyAPIEndpoint)

	// queue cache will be updated every 30 seconds
	resp, err := p.playerService.GetPlayerQueue(reqUrl, accessToken, 30)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}
