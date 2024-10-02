package service

import (
	"context"
	"fmt"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/CP-Payne/ecomstore/internal/models"
	"go.uber.org/zap"
)

type PaymentService struct {
	logger           *zap.Logger
	db               *database.Queries
	paymentProcessor models.PaymentProcessor
	orderSrv         *OrderService
	productSrv       *ProductService
	cartSrv          *CartService
}

func NewPaymentService(db *database.Queries, processor models.PaymentProcessor, orderSrv *OrderService, productSrv *ProductService, cartSrv *CartService) *PaymentService {
	return &PaymentService{
		logger:           config.GetLogger(),
		db:               db,
		paymentProcessor: processor,
		orderSrv:         orderSrv,
		productSrv:       productSrv,
		cartSrv:          cartSrv,
	}
}

func (p *PaymentService) CreateProcessorOrder(ctx context.Context, order *models.Order) (*models.OrderResult, error) {
	logger := p.logger.With(
		zap.String("method", "CreateProcessorOrder"),
		zap.String("orderID", order.ID.String()),
	)
	orderResult, err := p.paymentProcessor.CreateProcessorOrder(ctx, order)
	if err != nil {
		logger.Error("failed to create processor order")
		return nil, fmt.Errorf("failed to create processor order: %w", err)
	}

	err = p.orderSrv.UpdateOrderActionRequired(ctx, order.ID, orderResult.ID)
	if err != nil {
		logger.Error("failed to update order status and processor order ID")
		return nil, fmt.Errorf("failed to update order status and processor order ID: %w", err)
	}

	logger.Info("Succesfully created processor order", zap.String("ProcessorOrderID", orderResult.ID))
	return orderResult, nil
}

func (p *PaymentService) CaptureOrder(ctx context.Context, orderID string) error {

	logger := p.logger.With(
		zap.String("method", "CaptureOrder"),
		zap.String("orderID", orderID),
	)

	orderResult, err := p.paymentProcessor.CaptureOrder(ctx, orderID)
	if err != nil {
		logger.Error("failed to capture order", zap.Error(err))
		return fmt.Errorf("failed to capture order: %w", err)
	}

	err = p.orderSrv.UpdateOrderCompleted(ctx, orderResult)
	if err != nil {
		logger.Error("failed to update order status to completed", zap.Error(err))
		return fmt.Errorf("failed to update order status: %w", err)
	}

	order, err := p.orderSrv.GetOrderByProcessorOrderID(ctx, orderResult.ID)
	if err != nil {
		logger.Error("failed to get order by processor order ID", zap.Error(err))
		return fmt.Errorf("failed to retrieve order by processor order ID: %w", err)
	}

	orderItems := order.OrderItems

	// Update product Stock
	for _, item := range orderItems {
		err = p.productSrv.UpdateStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			logger.Error("failed to update product stock", zap.Error(err), zap.String("productID", item.ProductID.String()))
		}
	}

	if order.CartID != nil {
		err = p.cartSrv.DeleteCart(ctx, *order.CartID)
		if err != nil {
			logger.Info("failed to delete cart", zap.Error(err))
			return fmt.Errorf("failed to delete cart: %w", err)
		}
	}

	logger.Info("succesfully captured order")
	return nil
}
