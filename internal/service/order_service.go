package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/CP-Payne/ecomstore/internal/models"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

var globalShipping float32 = 0

type OrderService struct {
	logger *zap.Logger
	db     *database.Queries
	sqlDB  *sql.DB
}

func NewOrderService(db *database.Queries, sqlDB *sql.DB) *OrderService {
	return &OrderService{
		logger: config.GetLogger(),
		sqlDB:  sqlDB,
		db:     db,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, cart models.Cart) error {
	productTotal, err := s.getCartTotal(&cart)
	if err != nil {
		s.logger.Error("failed to calculate cart total", zap.Error(err))
		return fmt.Errorf("failed to calculate cart total: %w", err)
	}

	// Start transaction
	tx, err := s.sqlDB.Begin()
	if err != nil {
		s.logger.Error("failed to begin transaction", zap.Error(err))
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			s.logger.Error("failed to rollback transaction", zap.Error(err))
		}
	}()
	qtx := s.db.WithTx(tx)

	orderId, err := qtx.CreateOrder(ctx, database.CreateOrderParams{
		ID:            uuid.New(),
		UserID:        cart.UserID,
		ProductTotal:  floatToString(productTotal),
		Status:        "created",
		PaymentMethod: "paypal",
		ShippingPrice: floatToString(globalShipping),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	})
	if err != nil {
		s.logger.Error("failed to create order", zap.Error(err))
		return fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range cart.Items {
		if err := qtx.CreateOrderItem(ctx, database.CreateOrderItemParams{
			ID:        uuid.New(),
			OrderID:   orderId,
			ProductID: item.ProductID,
			Price:     floatToString(item.Price),
			Quantity:  int32(item.Quantity),
		}); err != nil {
			s.logger.Error("failed to create order item", zap.Error(err))
			return fmt.Errorf("failed to create order item for product %s: %w", item.ProductID, err)

		}
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("failed to commit transaction", zap.Error(err))
		return fmt.Errorf("failed to commit: %w", err)

	}

	return nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	return models.Order{}, nil
}

func floatToString(f float32) string {
	return fmt.Sprintf("%.2f", f)
}

func (p *OrderService) getCartTotal(cart *models.Cart) (float32, error) {
	if cart == nil {
		p.logger.Error("failed to calculate cart total. cart cannot be nil")
		return 0, fmt.Errorf("failed to calculate cart total -> cart is nil: %w", errors.New("nil cart"))
	}
	var cartTotal float32 = 0
	for _, ci := range cart.Items {
		cartTotal += ci.Price * float32(ci.Quantity)
	}

	return cartTotal, nil
}
