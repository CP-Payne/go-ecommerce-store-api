package models

import (
	"strconv"
	"time"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Product struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Price        float32   `json:"price"`
	Brand        string    `json:"brand"`
	Sku          string    `json:"sku"`
	Stock        int       `json:"stock"`
	CategoryID   uuid.UUID `json:"categoryId"`
	ImageURL     string    `json:"imageUrl"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	UpdatedAt    time.Time `json:"updatedAt,omitempty"`
}

// Database Product to product mappings
func DatabaseProductToProduct(product database.Product) Product {
	price, err := strconv.ParseFloat(product.Price, 32)
	if err != nil {
		config.GetLogger().Error("failed to parse string price to float", zap.Error(err))
		return Product{}
	}

	return Product{
		ID:           product.ID,
		Name:         product.Name,
		Description:  NullStringToString(product.Description),
		Price:        float32(price),
		Brand:        NullStringToString(product.Brand),
		Sku:          product.Sku,
		Stock:        int(product.StockQuantity),
		CategoryID:   product.ID,
		ImageURL:     NullStringToString(product.ImageUrl),
		ThumbnailURL: NullStringToString(product.ThumbnailUrl),
		IsActive:     product.IsActive,
		CreatedAt:    product.CreatedAt,
		UpdatedAt:    product.UpdatedAt,
	}
}
