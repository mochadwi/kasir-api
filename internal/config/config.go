package config

import (
	"fmt"

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

	viper.SetConfigFile(".env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("..")
	viper.AddConfigPath("../..")

	viper.AutomaticEnv()

	// Try to read .env file, but don't fail if it doesn't exist
	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
