package service

import (
	"context"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/CP-Payne/ecomstore/internal/models"
	"go.uber.org/zap"
)

const shippingPrice = 0.00

type PaymentService struct {
	logger           *zap.Logger
	db               *database.Queries
	paymentProcessor models.PaymentProcessor
}

func NewPaymentService(db *database.Queries, processor models.PaymentProcessor) *PaymentService {
	return &PaymentService{
		logger:           config.GetLogger(),
		db:               db,
		paymentProcessor: processor,
	}
}

func (p *PaymentService) CreateProductOrder(ctx context.Context, product *models.Product, quantity int) (*models.OrderResult, error) {
	orderResult, err := p.paymentProcessor.CreateProductOrder(ctx, product, quantity, shippingPrice)
	if err != nil {
		return nil, err
	}
	return orderResult, nil
}

func (p *PaymentService) CreateCartOrder(ctx context.Context, cart *models.Cart) (*models.OrderResult, error) {
	orderResult, err := p.paymentProcessor.CreateCartOrder(ctx, cart, shippingPrice)
	if err != nil {
		return nil, err
	}
	return orderResult, nil
}

func (p *PaymentService) CaptureOrder(ctx context.Context, orderID string) error {
	err := p.paymentProcessor.CaptureOrder(ctx, orderID)
	if err != nil {
		return err
	}
	return nil
}

func (p *PaymentService) GetOrderDetails(ctx context.Context, orderID string) error {
	err := p.paymentProcessor.GetOrderDetails(ctx, orderID)
	if err != nil {
		return err
	}
	return nil
}
