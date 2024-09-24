package models

import (
	"time"

	"github.com/google/uuid"
)

// type Order struct {
// 	ID            uuid.UUID  `json:"id"`
// 	UserID        uuid.UUID  `json:"userId"`
// 	Items         []CartItem `json:"items"`
// 	TotalAmount   float32    `json:"totalAmount"`
// 	Status        string     `json:"status"`
// 	PaymentMethod string     `json:"paymentMethod"`
// 	PaymentID     string     `json:"paymentId"`
// 	CreatedAt     time.Time  `json:"createdAt"`
// 	UpdatedAt     time.Time  `json:"updatedAt"`
// }

type Order struct {
	ID               uuid.UUID   `json:"id"`
	ProductTotal     float32     `json:"productTotal"`
	OrderTotal       float32     `json:"orderTotal"`
	ProcessorOrderID string      `json:"processorOrderId"`
	Status           string      `json:"status"`
	UserID           uuid.UUID   `json:"userId"`
	OrderItems       []OrderItem `json:"items"`
	PaymentEmail     string      `json:"paymentEmail"`
	PaymentMethod    string      `json:"paymentMethod"`
	PayerID          string      `json:"payerId"`
	ShippingPrice    float32     `json:"shippingPrice"`
	CartID           *uuid.UUID  `json:"cartId;omitempty"`
	CreatedAt        time.Time   `json:"createdAt"`
	UpdatedAt        time.Time   `json:"updatedAt"`
}

// TODO: Need to set PayerID, PaymentEmail, ProcessorOrderID

type OrderItem struct {
	ProductID uuid.UUID `json:"productId"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	Price     float32   `json:"price"`
}
