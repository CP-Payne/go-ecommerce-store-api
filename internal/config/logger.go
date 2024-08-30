package config

import (
	"sync"

	"go.uber.org/zap"
)

var (
	once     sync.Once
	instance *zap.Logger
)

// GetLogger returns a singleton instance of the logger
func GetLogger() *zap.Logger {
	once.Do(func() {
		// Initialise the logger only once
		logger := zap.Must(zap.NewProduction())
		instance = logger
	})

	return instance
}
