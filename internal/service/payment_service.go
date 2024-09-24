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
	// orderResult, err := p.paymentProcessor.CreateCartOrder(ctx, cart, shippingPrice)
	orderResult, err := p.paymentProcessor.CreateProcessorOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	err = p.orderSrv.UpdateOrderActionRequired(ctx, order.ID, orderResult.ID)
	if err != nil {
		return nil, err
	}
	return orderResult, nil
}

func (p *PaymentService) CaptureOrder(ctx context.Context, orderID string) error {
	orderResult, err := p.paymentProcessor.CaptureOrder(ctx, orderID)
	if err != nil {
		return err
	}

	err = p.orderSrv.UpdateOrderCompleted(ctx, orderResult)
	if err != nil {
		return err
	}

	order, err := p.orderSrv.GetOrderByProcessorOrderID(ctx, orderResult.ID)
	if err != nil {
		return err
	}

	orderItems := order.OrderItems

	// Update product Stock
	for _, item := range orderItems {
		err = p.productSrv.UpdateStock(ctx, item.ProductID, item.Quantity)
		if err != nil {
			p.logger.Error("failed to update product stock", zap.Error(err), zap.String("productID", item.ProductID.String()))
		}
	}

	if order.CartID != nil {
		err = p.cartSrv.DeleteCart(ctx, *order.CartID)
		if err != nil {
			return err
		}
	}

	return nil
}

// func (p *PaymentService) GetOrderDetails(ctx context.Context, orderID string) error {
// 	err := p.paymentProcessor.GetOrderDetails(ctx, orderID)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
