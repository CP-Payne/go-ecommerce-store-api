package service

import (
	"context"
	"fmt"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/CP-Payne/ecomstore/internal/models"
	"github.com/CP-Payne/ecomstore/internal/utils/apperrors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProductService struct {
	logger *zap.Logger
	db     *database.Queries
}

func NewProductService(db *database.Queries) *ProductService {
	return &ProductService{
		logger: config.GetLogger(),
		db:     db,
	}
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (models.Product, error) {
	product, err := s.db.GetProduct(ctx, id)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			return models.Product{}, apperrors.ErrNotFound
		}
		s.logger.Error("failed to retrieve product from db", zap.Error(err), zap.String("productID", id.String()))
		return models.Product{}, fmt.Errorf("failed to retrieve product: %w", err)

	}
	return models.DatabaseProductToProduct(product), nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	products, err := s.db.GetAllProducts(ctx)
	if err != nil {
		s.logger.Error("failed to retrieve products from db", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve and process products: %w", err)
	}

	return models.DatabaseProductsToProducts(products), nil
}
