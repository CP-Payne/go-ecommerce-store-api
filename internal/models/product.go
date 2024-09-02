package models

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Product struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	Price          float32         `json:"price"`
	Brand          string          `json:"brand"`
	Sku            string          `json:"sku"`
	Stock          int             `json:"stock"`
	CategoryID     uuid.UUID       `json:"categoryId"`
	ImageURL       string          `json:"imageUrl"`
	ThumbnailURL   string          `json:"thumbnailUrl"`
	Specifications json.RawMessage `json:"specifications"`
	Variants       json.RawMessage `json:"variants"`
}

type ProductWithMetadata struct {
	Product
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Database Product to product mappings
func DatabaseProductToProduct(product database.Product, includeMetadata bool) interface{} {
	price, err := strconv.ParseFloat(product.Price, 32)
	if err != nil {
		config.GetLogger().Error("failed to parse string price to float", zap.Error(err))
		return Product{}
	}

	baseProduct := Product{
		ID:             product.ID,
		Name:           product.Name,
		Description:    NullStringToString(product.Description),
		Price:          float32(price),
		Brand:          NullStringToString(product.Brand),
		Sku:            product.Sku,
		Stock:          int(product.StockQuantity),
		CategoryID:     product.ID,
		ImageURL:       NullStringToString(product.ImageUrl),
		ThumbnailURL:   NullStringToString(product.ThumbnailUrl),
		Specifications: NullRawMessageToRawMessage(product.Specifications),
		Variants:       NullRawMessageToRawMessage(product.Variants),
		// IsActive:       product.IsActive,
		// CreatedAt:      product.CreatedAt,
		// UpdatedAt:      product.UpdatedAt,
	}

	if includeMetadata {
		return ProductWithMetadata{
			Product:   baseProduct,
			IsActive:  product.IsActive,
			CreatedAt: product.CreatedAt,
			UpdatedAt: product.UpdatedAt,
		}
	}

	return baseProduct
}

func DatabaseProductsToProducts(dbProducts []database.Product, includeMetadata bool) interface{} {
	if len(dbProducts) == 0 {
		if includeMetadata {
			return []ProductWithMetadata{}
		}
		return []Product{}
	}
	if includeMetadata {
		products := make([]ProductWithMetadata, 0, len(dbProducts))
		for _, dbProd := range dbProducts {
			product := DatabaseProductToProduct(dbProd, includeMetadata).(ProductWithMetadata)
			products = append(products, product)
		}
		return products

	} else {
		products := make([]Product, 0, len(dbProducts))
		for _, dbProd := range dbProducts {
			product := DatabaseProductToProduct(dbProd, includeMetadata).(Product)
			products = append(products, product)
		}
		return products

	}
}
