package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Config struct {
	Port               string   `env:"PORT" envDefault:"8080"`
	GinMode            string   `env:"GIN_MODE" envDefault:"debug"`
	LogLevel           string   `env:"LOG_LEVEL" envDefault:"info"`
	TrustedProxies     []string `env:"TRUSTED_PROXIES" envDefault:"10.0.0.0/8,172.16.0.0/12,192.168.0.0/16" envSeparator:","`
	ClientId           string   `env:"SPOTIFY_CLIENT_ID,required"`
	ClientSecret       string   `env:"SPOTIFY_CLIENT_SECRET,required"`
	RedirectUri        string   `env:"SPOTIFY_REDIRECT_URI,required"`
	State              string   `env:"SPOTIFY_STATE"`
	Scope              []string `env:"SPOTIFY_SCOPE" envSeparator:","`
	RefreshToken       string   `env:"SPOTIFY_REFRESH_TOKEN"`
	RefreshTokenOutput bool     `env:"SPOTIFY_REFRESH_TOKEN_OUTPUT" envDefault:"false"`
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
