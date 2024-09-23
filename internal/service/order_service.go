package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/CP-Payne/ecomstore/internal/models"
	"github.com/CP-Payne/ecomstore/internal/utils/apperrors"
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

func (s *OrderService) CreateOrder(ctx context.Context, cart models.Cart, tempCart bool) (models.Order, error) {
	productTotal, err := s.getCartTotal(&cart)
	if err != nil {
		s.logger.Error("failed to calculate cart total", zap.Error(err))
		return models.Order{}, fmt.Errorf("failed to calculate cart total: %w", err)
	}

	orderTotal := productTotal + globalShipping

	cartID := uuid.NullUUID{
		Valid: true,
		UUID:  cart.ID,
	}

	if tempCart {
		cartID.Valid = false
	}

	// Start transaction
	tx, err := s.sqlDB.Begin()
	if err != nil {
		s.logger.Error("failed to begin transaction", zap.Error(err))
		return models.Order{}, fmt.Errorf("failed to begin transaction: %w", err)
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
		OrderTotal:    floatToString(orderTotal),
		Status:        "created",
		PaymentMethod: "paypal",
		ShippingPrice: floatToString(globalShipping),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		CartID:        cartID,
	})
	if err != nil {
		s.logger.Error("failed to create order", zap.Error(err))
		return models.Order{}, fmt.Errorf("failed to create order: %w", err)
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
			return models.Order{}, fmt.Errorf("failed to create order item for product %s: %w", item.ProductID, err)

		}
	}

	if err := tx.Commit(); err != nil {
		s.logger.Error("failed to commit transaction", zap.Error(err))
		return models.Order{}, fmt.Errorf("failed to commit: %w", err)

	}

	// Get newly created order
	order, err := s.GetOrderByID(ctx, orderId)
	if err != nil {
		s.logger.Error("failed to retrieve order", zap.Error(err))
		return models.Order{}, fmt.Errorf("failed to retrieve new order: %w", err)

	}

	return order, nil
}

func (s *OrderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (models.Order, error) {
	orderRecord, err := s.db.GetOrderByID(ctx, orderID)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			return models.Order{}, apperrors.ErrNotFound
		}
		s.logger.Error("failed to retrieve order from db", zap.Error(err), zap.String("orderID", orderID.String()))
		return models.Order{}, fmt.Errorf("failed to retrieve order from db: %w", err)
	}

	return s.DatabaseOrderToOrder(ctx, orderRecord)
}

func (s *OrderService) GetOrderByProcessorOrderID(ctx context.Context, processorOrderID string) (models.Order, error) {
	orderRecord, err := s.db.GetOrderByProcessorOrderID(ctx, sql.NullString{
		Valid:  true,
		String: processorOrderID,
	})
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			return models.Order{}, apperrors.ErrNotFound
		}
		s.logger.Error("failed to retrieve order from db", zap.Error(err), zap.String("processorOrderID", processorOrderID))
		return models.Order{}, fmt.Errorf("failed to retrieve order from db: %w", err)
	}

	return s.DatabaseOrderToOrder(ctx, orderRecord)
}

// Process database Order
func (s *OrderService) DatabaseOrderToOrder(ctx context.Context, orderRecord database.Order) (models.Order, error) {
	orderItemsRecord, err := s.db.GetOrderItemsByOrderID(ctx, orderRecord.ID)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			return models.Order{}, apperrors.ErrNotFound
		}
		s.logger.Error("failed to retrieve order items from db", zap.Error(err), zap.String("orderID", orderRecord.ID.String()))
		return models.Order{}, fmt.Errorf("failed to retrieve order items from db: %w", err)
	}

	productTotal, err := stringToFloat32(orderRecord.ProductTotal)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to get product total: %w", err)
	}

	shippingPrice, err := stringToFloat32(orderRecord.ShippingPrice)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to get shipping price: %w", err)
	}

	orderTotal, err := stringToFloat32(orderRecord.OrderTotal)
	if err != nil {
		return models.Order{}, fmt.Errorf("failed to get order total: %w", err)
	}

	// Create order
	order := models.Order{
		ID:               orderRecord.ID,
		ProductTotal:     productTotal,
		OrderTotal:       orderTotal,
		Status:           orderRecord.Status,
		UserID:           orderRecord.UserID,
		PaymentMethod:    orderRecord.PaymentMethod,
		ProcessorOrderID: sqlNullStringToString(orderRecord.ProcessorOrderID),
		PaymentEmail:     sqlNullStringToString(orderRecord.PaymentEmail),
		PayerID:          sqlNullStringToString(orderRecord.PayerID),
		ShippingPrice:    shippingPrice,
		CartID:           nullUuidToUuid(orderRecord.CartID),
		CreatedAt:        orderRecord.CreatedAt,
		UpdatedAt:        orderRecord.UpdatedAt,
	}

	orderItems := make([]models.OrderItem, 0, len(orderItemsRecord))
	for _, item := range orderItemsRecord {
		priceF, err := stringToFloat32(item.Price)
		if err != nil {
			return models.Order{}, fmt.Errorf("failed to get item price: %w", err)
		}
		orderItem := models.OrderItem{
			ProductID: item.ProductID,
			Name:      item.Name,
			Quantity:  int(item.Quantity),
			Price:     priceF,
		}
		orderItems = append(orderItems, orderItem)
	}

	order.OrderItems = orderItems

	return order, nil
}

func (s *OrderService) UpdateOrderActionRequired(ctx context.Context, orderID uuid.UUID, processorOrderID string) error {
	err := s.db.SetProcessorIDAndStatus(ctx, database.SetProcessorIDAndStatusParams{
		ProcessorOrderID: sql.NullString{
			Valid:  true,
			String: processorOrderID,
		},
		Status:    "PAYER_ACTION_REQUIRED",
		ID:        orderID,
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to update order: %w", err)
	}
	return nil
}

func (s *OrderService) UpdateOrderCompleted(ctx context.Context, orderResult *models.OrderResult) error {
	err := s.db.SetOrderCompleted(ctx, database.SetOrderCompletedParams{
		Status: "COMPLETED",
		PaymentEmail: sql.NullString{
			Valid:  true,
			String: orderResult.PaymentEmail,
		},
		PayerID: sql.NullString{
			Valid:  true,
			String: orderResult.PayerID,
		},
		ProcessorOrderID: sql.NullString{
			Valid:  true,
			String: orderResult.ID,
		},
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to update order status to complete: %w", err)
	}
	return nil
}

func stringToFloat32(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return 0, err
	}
	return float32(f), nil
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
