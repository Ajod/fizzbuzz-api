package config

import (
	"fizzbuzz-api/internal/fizzbuzzapi/logger"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Port string `envconfig:"PORT" default:"4255"`
	Host string `envconfig:"HOST" default:"localhost"`

	MaxFizzBuzzLimit int    `envconfig:"MAX_FIZZBUZZ_LIMIT" default:"100000"` // Max limit for FizzBuzz generation
	MaxStringLength  int    `envconfig:"MAX_STRING_LENGTH" default:"30"`      // Max length for Str1 and Str2
	StatsStorage     string `envconfig:"STATS_STORAGE" default:"inmemory"`    // Storage type for stats: "inmemory" or "file"
}

func LoadConfig(log logger.Logger) (*Config, error) {
	cfg := Config{}

	// Load environment variables into the config struct
	// Environment variables must be prefixed with "FBAPI_"
	err := envconfig.Process("FBAPI", &cfg)
	if err != nil {
		return nil, err
	}

	log.Info("config loaded", "config", cfg)
	return &cfg, nil
}
