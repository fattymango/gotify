package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

const (
	DB_PATH               = "DB_PATH"
	LOG_FILE              = "LOG_FILE"
	SPOTIFY_CLIENT_ID     = "SPOTIFY_CLIENT_ID"
	SPOTIFY_CLIENT_SECRET = "SPOTIFY_CLIENT_SECRET"
	SPOTIFY_REDIRECT_URL  = "SPOTIFY_REDIRECT_URL"
	SPOTIFY_STATE         = "SPOTIFY_STATE"
	SERVER_ADDRESS        = "SERVER_ADDRESS"
)

// Config for the application
type Config struct {
	Logger  Logger
	DB      DB
	Setting Setting
	Spotify Spotify
	Server  Server
}

type Logger struct {
	DisableCaller     bool   `envconfig:"DISABLE_CALLER"`
	DisableStacktrace bool   `envconfig:"DISABLE_STACKTRACE"`
	Encoding          string `envconfig:"ENCODING"`
	Level             string `envconfig:"LEVEL"`
	LogFile           string `envconfig:"LOG_FILE" required:"true"`
}

type DB struct {
	Path string `envconfig:"DB_PATH" required:"true"`
}

type Spotify struct {
	ClientID     string `envconfig:"SPOTIFY_CLIENT_ID" required:"true"`
	ClientSecret string `envconfig:"SPOTIFY_CLIENT_SECRET" required:"true"`
	RedirectURL  string `envconfig:"SPOTIFY_REDIRECT_URL"`
	State        string `envconfig:"SPOTIFY_STATE"`
}

type Setting struct {
	LocalPath string `envconfig:"LOCAL_PATH" default:"./locales/*/*"`
	Version   string `envconfig:"VERSION" default:"1.0.0"`
	DEBUG     bool   `envconfig:"DEBUG" default:"false"`
}

type Server struct {
	Address string `envconfig:"SERVER_ADDRESS" required:"true"`
}

// NewConfig loads configuration from environment variables
func NewConfig() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c) // Reads environment variables and populates the struct
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &c, nil
}
