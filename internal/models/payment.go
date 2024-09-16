package models

import (
	"context"
)

type PaymentProcessor interface {
	// CreateCartOrder()
	CreateProductOrder(ctx context.Context, product *Product, quantity int, shippingPrice float32) (*OrderResult, error)
	CreateCartOrder(ctx context.Context, cart *Cart, shippingPrice float32) (*OrderResult, error)
	CaptureOrder(ctx context.Context, orderID string) error
	// CompleteOrder(orderId int)
}

type OrderResult struct {
	ID          string
	ApproveLink string
	Status      string
}
