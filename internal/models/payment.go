package models

import (
	"context"
)

type PaymentProcessor interface {
	CaptureOrder(ctx context.Context, orderID string) (*OrderResult, error)
	// GetOrderDetails(ctx context.Context, orderID string) error
	CreateProcessorOrder(ctx context.Context, order *Order) (*OrderResult, error)
}

type OrderResult struct {
	ID           string
	ApproveLink  string
	Status       string
	PaymentEmail string
	PayerID      string
}
