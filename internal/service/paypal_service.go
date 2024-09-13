package service

import (
	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/plutov/paypal/v4"
	"go.uber.org/zap"
)

type PaypalService struct {
	logger *zap.Logger
	db     *database.Queries
	client *paypal.Client
}

func NewPaypalService(db *database.Queries, client *paypal.Client) *PaypalService {
	return &PaypalService{
		logger: config.GetLogger(),
		db:     db,
		client: client,
	}
}
