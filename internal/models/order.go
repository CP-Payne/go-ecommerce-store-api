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
	ID            uuid.UUID `json:"id"`
	TotalAmount   float32   `json:"totalAmount"`
	Status        string    `json:"status"`
	UserID        uuid.UUID `json:"userId"`
	OrderItems    []OrderItem
	PaymentEmail  string
	PaymentMethod string `json:"paymentMethod"`
	PayerID       string
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	ShippingPrice float32
}

type OrderItem struct {
	ProductID uuid.UUID
	Name      string
	Quantity  int
	Price     float32
}
