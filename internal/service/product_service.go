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

func (s *ProductService) ProductExists(ctx context.Context, id uuid.UUID) (bool, error) {
	logger := s.logger.With(
		zap.String("method", "ProductExists"),
		zap.String("productID", id.String()),
	)
	exists, err := s.db.ProductExists(ctx, id)
	if err != nil {
		logger.Error("failed to check if product exists", zap.Error(err), zap.String("productID", id.String()))
		return false, fmt.Errorf("failed to check if product exists: %w", err)
	}

	return exists, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id uuid.UUID) (models.Product, error) {
	logger := s.logger.With(
		zap.String("method", "GetProduct"),
		zap.String("productID", id.String()),
	)
	product, err := s.db.GetProduct(ctx, id)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			logger.Info("user attempted to retrieve product that does not exist")
			return models.Product{}, fmt.Errorf("failed to retrieve product: %w", apperrors.ErrNotFound)
		}
		logger.Error("failed to retrieve product", zap.Error(err))
		return models.Product{}, fmt.Errorf("failed to retrieve product: %w", err)

	}

	p, ok := models.DatabaseProductToProduct(product, false).(models.Product)
	if !ok {
		logger.Error("failed to convert database product into product model")
		return models.Product{}, fmt.Errorf("failed to process database product: %w", err)
	}

	return p, nil
}

func (s *ProductService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	logger := s.logger.With(
		zap.String("method", "GetAllProducts"),
	)

	products, err := s.db.GetAllProducts(ctx)
	if err != nil {
		logger.Error("failed to retrieve products", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve products: %w", err)
	}

	pl, ok := models.DatabaseProductsToProducts(products, false).([]models.Product)
	if !ok {
		logger.Error("failed to convert database products into product model")
		return []models.Product{}, fmt.Errorf("failed to process database products: %w", err)
	}

	return pl, nil
}

func (s *ProductService) GetProductCategories(ctx context.Context) ([]models.Category, error) {
	logger := s.logger.With(
		zap.String("method", "GetProductCategories"),
	)
	categories, err := s.db.GetProductCategories(ctx)
	if err != nil {
		logger.Error("failed to retrieve product categories", zap.Error(err))
		return nil, fmt.Errorf("failed to retrieve product categories: %w", err)
	}
	return models.DatabaseCategoriesToCategories(categories), nil
}

func (s *ProductService) GetProductsByCategory(ctx context.Context, categoryID uuid.UUID) ([]models.Product, error) {
	logger := s.logger.With(
		zap.String("method", "GetProductsByCategory"),
		zap.String("categoryID", categoryID.String()),
	)

	products, err := s.db.GetProductsByCategory(ctx, categoryID)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			logger.Info("attempted to find products for category that does not exist", zap.Error(err))
			return []models.Product{}, fmt.Errorf("failed to retrieve products by category: %w", apperrors.ErrNotFound)
		}

		logger.Error("failed to retrieve products by category", zap.Error(err))
		return []models.Product{}, fmt.Errorf("failed to retrieve products by category: %w", err)
	}

	pl, ok := models.DatabaseProductsToProducts(products, false).([]models.Product)
	if !ok {

		logger.Error("failed to convert database products into products model")
		return []models.Product{}, fmt.Errorf("failed to process database products: %w", err)
	}
	return pl, nil
}

func (s *ProductService) UpdateStock(ctx context.Context, productID uuid.UUID, reduceBy int) error {
	logger := s.logger.With(
		zap.String("method", "UpdateStock"),
		zap.String("productID", productID.String()),
		zap.Int("amount", reduceBy),
	)
	err := s.db.UpdateStock(ctx, database.UpdateStockParams{
		ID:            productID,
		StockQuantity: int32(reduceBy),
	})
	if err != nil {
		logger.Error("failed to update product stock")
		return fmt.Errorf("failed to update product's stock: %w", err)
	}

	logger.Info("product stock successfully updated")
	return nil
}
