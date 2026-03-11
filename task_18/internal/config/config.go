package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "EVENT_APP"

type Config struct {
	Port int `envconfig:"PORT" required:"true"`
}

func Load() (*Config, error) {
	loadEnvFile()

	var cfg Config

	err := envconfig.Process(envPrefix, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func loadEnvFile() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(".env"); err != nil {
			log.Printf("Failed to load .env: %v", err)
		}
	}
}
