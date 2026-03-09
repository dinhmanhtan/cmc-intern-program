package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_DSN string
	PORT   string
}

func Load() *Config {

	godotenv.Load()

	return &Config{
		DB_DSN: os.Getenv("DB_DSN"),
		PORT:   os.Getenv("PORT"),
	}
}
