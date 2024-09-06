package service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/CP-Payne/ecomstore/internal/models"
	"github.com/CP-Payne/ecomstore/internal/utils/apperrors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ReviewService struct {
	logger *zap.Logger
	db     *database.Queries
}

func NewReviewService(db *database.Queries) *ReviewService {
	return &ReviewService{
		logger: config.GetLogger(),
		db:     db,
	}
}

func (s *ReviewService) PostReview(ctx context.Context, title, reviewText string, rating int, productID, userID uuid.UUID) (models.Review, error) {
	// TODO: make sure that the user can only add review to product he purchased
	// TODO: user should only be able to add a review if he has not yet reviewed the product or the review has a status of deleted

	// Check if user has already reviewed product
	reviewed, err := s.db.HasUserReviewedProduct(ctx, database.HasUserReviewedProductParams{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		s.logger.Error("failed to check if user already reviewed product", zap.Error(err), zap.String("userID", userID.String()),
			zap.String("productID", productID.String()))
		return models.Review{}, fmt.Errorf("failed to determine if user already reviewed product: %w", apperrors.ErrInternal)
	}

	if reviewed {
		return models.Review{}, fmt.Errorf("user already reviewed product: %w", apperrors.ErrConflict)
	}
	dbReview, err := s.db.InsertReview(ctx, database.InsertReviewParams{
		ID:         uuid.New(),
		Title:      sql.NullString{String: title, Valid: true},
		ReviewText: sql.NullString{String: reviewText, Valid: true},
		Rating:     int32(rating),
		ProductID:  productID,
		UserID:     userID,
		Deleted:    false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		s.logger.Error("failed to add product review", zap.Error(err), zap.String("UserID", userID.String()), zap.String("ProductID", productID.String()))
		return models.Review{}, fmt.Errorf("failed to add product review to db: %w", err)
	}
	review := models.DatabaseReviewToReview(dbReview)
	return review, nil
}

func (s *ReviewService) GetProductReviews(ctx context.Context, productID uuid.UUID) ([]models.Review, error) {
	dbReviews, err := s.db.GetProductReviews(ctx, productID)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			return []models.Review{}, nil
		}
		s.logger.Error("failed to retrieve product reviews from db", zap.Error(err), zap.String("productID", productID.String()))
		return []models.Review{}, fmt.Errorf("failed to retrieve product reviews: %w", err)
	}

	return models.DatabaseReviewsToReviews(dbReviews), nil
}

func (s *ReviewService) GetReviewByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) (models.Review, error) {
	dbReview, err := s.db.GetReviewByUserAndProduct(ctx, database.GetReviewByUserAndProductParams{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			return models.Review{}, fmt.Errorf("no user review for product found: %w", apperrors.ErrNotFound)
		}
		s.logger.Error("failed to retrieve user review for product", zap.Error(err),
			zap.String("userID", userID.String()), zap.String("productID", productID.String()))
		return models.Review{}, fmt.Errorf("failed to retrieve user review for product: %w", err)
	}
	return models.DatabaseReviewToReview(dbReview), nil
}
