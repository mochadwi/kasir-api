package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server struct {
		Port string `mapstructure:"PORT"`
	}
	Database struct {
		URL string `mapstructure:"DATABASE_URL"`
	}
}

// Load reads configuration from environment variables and .env file
func Load() (*Config, error) {
	viper.SetDefault("PORT", "8080")

	// First try to load from .env file
	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./cmd/api")

	// Enable automatic environment variable reading
	viper.AutomaticEnv()

	// Try to read .env file, but don't fail if it doesn't exist
	if err := viper.ReadInConfig(); err != nil {
		// .env file not found - that's okay, we'll use env vars
		fmt.Println("Note: .env file not found, using environment variables")
	} else {
		fmt.Println("Loaded config from:", viper.ConfigFileUsed())
	}

	var cfg Config

	// Get values directly from viper (handles both .env and env vars)
	cfg.Server.Port = viper.GetString("PORT")
	cfg.Database.URL = viper.GetString("DATABASE_URL")

	// Fallback to direct env var if still empty
	if cfg.Database.URL == "" {
		cfg.Database.URL = os.Getenv("DATABASE_URL")
	}
	if cfg.Server.Port == "" {
		cfg.Server.Port = os.Getenv("PORT")
	}

	// Clean up quotes from DATABASE_URL if present (Viper includes them if .env has quotes)
	cfg.Database.URL = strings.Trim(cfg.Database.URL, `"'`)
	cfg.Database.URL = strings.TrimSpace(cfg.Database.URL)

	return &cfg, nil
}
