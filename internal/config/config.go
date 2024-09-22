package config

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Config struct {
	Logger           *zap.Logger
	Port             string
	DB               *database.Queries
	SqlDB            *sql.DB
	PaymentProcessor *ProcessorConfig
}

type ProcessorConfig struct {
	ClientID     string
	ClientSecret string
	ReturnUrl    string
	CancelUrl    string
}

func New() *Config {
	logger := GetLogger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal("failed to load environment variables", zap.Error(err))
	}

	port := os.Getenv("PORT")

	// Database initialisation

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassord := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbName := os.Getenv("POSTGRES_DB")

	connString := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPassord, dbHost, dbName)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		logger.Fatal("failed to open database connection", zap.Error(err))
	}

	ppClientID := os.Getenv("PAYPAL_CLIENT")
	ppClientSecret := os.Getenv("PAYPAL_SECRET")
	paypalReturnUrl := os.Getenv("PAYPAL_RETURN_URL")
	paypalCancelUrl := os.Getenv("PAYPAL_CANCEL_URL")

	return &Config{
		Port:   port,
		Logger: logger,
		SqlDB:  db,
		DB:     database.New(db),
		PaymentProcessor: &ProcessorConfig{
			ClientID:     ppClientID,
			ClientSecret: ppClientSecret,
			ReturnUrl:    paypalReturnUrl,
			CancelUrl:    paypalCancelUrl,
		},
	}
}
