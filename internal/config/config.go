package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

type Config struct {
	Logger *zap.Logger
	Port   string
}

func New() *Config {
	logger := GetLogger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal("failed to load environment variables", zap.Error(err))
	}

	port := os.Getenv("PORT")

	return &Config{
		Port:   port,
		Logger: logger,
	}
}
