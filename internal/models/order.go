package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID  `json:"id"`
	UserID        uuid.UUID  `json:"userId"`
	Items         []CartItem `json:"items"`
	TotalAmount   float32    `json:"totalAmount"`
	Status        string     `json:"status"`
	PaymentMethod string     `json:"paymentMethod"`
	PaymentID     string     `json:"paymentId"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}
