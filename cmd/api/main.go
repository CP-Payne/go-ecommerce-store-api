package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/CP-Payne/ecomstore/internal/api"
	"github.com/CP-Payne/ecomstore/internal/config"
	"go.uber.org/zap"
)

func main() {
	cfg := config.New()

	killSignal := make(chan os.Signal, 1)

	signal.Notify(killSignal, os.Interrupt, syscall.SIGTERM)

	r := api.SetupRouter(cfg.DB)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	// Start server in goroutine
	go func() {
		err := server.ListenAndServe()

		if errors.Is(err, http.ErrServerClosed) {
			cfg.Logger.Info("Server shutdown complete")
			// log.Printf("Server shutdown complete")
		} else if err != nil {
			cfg.Logger.Fatal("server failed to start", zap.Error(err))
		}
	}()

	cfg.Logger.Info("Server started...")

	// Wait for killsignal
	<-killSignal

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		cfg.Logger.Fatal("server shutdown failed", zap.Error(err))
	}
}
