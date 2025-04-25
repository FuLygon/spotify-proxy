package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	AccessPort         string   `env:"ACCESS_PORT" envDefault:"8000"`
	ProxyPort          string   `env:"PROXY_PORT" envDefault:"8001"`
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

var validate = validator.New()

type ProxyRouteConfig struct {
	Path    string   `yaml:"path" validate:"required"`
	Methods []string `yaml:"methods" validate:"min=1,dive,oneof=GET POST PUT DELETE PATCH OPTIONS"`
}

type ProxyRoutesConfig struct {
	ProxyRoutes []ProxyRouteConfig `yaml:"routes" validate:"dive"`
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

func LoadProxyRoutesConfig(path string) (*ProxyRoutesConfig, error) {
	// Skip if config file does not exist
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config ProxyRoutesConfig
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	err = validate.Struct(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
