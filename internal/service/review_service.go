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

func (s *ReviewService) PostReview(ctx context.Context, title, reviewText string, rating int, anonymous bool, productID, userID uuid.UUID) (models.Review, error) {
	logger := s.logger.With(
		zap.String("method", "PostReview"),
		zap.String("userID", userID.String()),
		zap.String("productID", productID.String()),
	)

	reviewed, err := s.db.HasUserReviewedProduct(ctx, database.HasUserReviewedProductParams{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		logger.Error("failed to check if user already reviewed product", zap.Error(err))
		return models.Review{}, fmt.Errorf("failed to check existing review: %w", apperrors.ErrInternal)
	}

	if reviewed {
		logger.Info("user attempted to review product again")
		return models.Review{}, apperrors.ErrConflict
	}
	dbReview, err := s.db.InsertReview(ctx, database.InsertReviewParams{
		ID:         uuid.New(),
		Title:      sql.NullString{String: title, Valid: true},
		ReviewText: sql.NullString{String: reviewText, Valid: true},
		Rating:     int32(rating),
		ProductID:  productID,
		UserID:     userID,
		Deleted:    false,
		Anonymous:  anonymous,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})
	if err != nil {
		logger.Error("failed to insert review", zap.Error(err))
		return models.Review{}, fmt.Errorf("failed to add review: %w", err)
	}
	review := models.DatabaseReviewToReview(dbReview)
	logger.Info("review added successfully")
	return review, nil
}

func (s *ReviewService) GetProductReviews(ctx context.Context, productID uuid.UUID) ([]models.ReviewDisplay, error) {
	logger := s.logger.With(
		zap.String("method", "GetProductReviews"),
		zap.String("productID", productID.String()),
	)
	dbReviews, err := s.db.GetProductReviews(ctx, productID)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			logger.Info("user attempted to get reviews for a product that does not exist")
			return []models.ReviewDisplay{}, nil
		}
		logger.Error("failed to retrieve product reviews", zap.Error(err))
		return []models.ReviewDisplay{}, fmt.Errorf("failed to retrieve product reviews: %w", err)
	}

	return models.DatabaseProductReviewsToReviewDisplays(dbReviews), nil
}

func (s *ReviewService) GetReviewByUserAndProduct(ctx context.Context, userID, productID uuid.UUID) (models.ReviewDisplay, error) {

	logger := s.logger.With(
		zap.String("method", "GetReviewByUserAndProduct"),
		zap.String("userID", userID.String()),
		zap.String("productID", productID.String()),
	)

	dbReview, err := s.db.GetReviewByUserAndProduct(ctx, database.GetReviewByUserAndProductParams{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			logger.Error("user review for product not found", zap.Error(err))
			return models.ReviewDisplay{}, fmt.Errorf("user review for product not found: %w", apperrors.ErrNotFound)
		}
		logger.Error("failed to retrieve user review", zap.Error(err))
		// TODO: update all errors and return apperrors
		return models.ReviewDisplay{}, fmt.Errorf("failed to retrieve user review: %w", apperrors.ErrInternal)
	}
	return models.DatabaseUserProductReviewToReviewDisplay(dbReview), nil
}

func (s *ReviewService) DeleteReview(ctx context.Context, userID, productID uuid.UUID) error {
	logger := s.logger.With(
		zap.String("method", "DeleteReview"),
		zap.String("userID", userID.String()),
		zap.String("productID", productID.String()),
	)
	err := s.db.SetReviewStatusDeleted(ctx, database.SetReviewStatusDeletedParams{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		logger.Error("failed to update review status to deleted", zap.Error(err))
		return fmt.Errorf("failed to delete review: %w", apperrors.ErrInternal)
	}
	logger.Info("review successfully deleted")
	return nil
}

func (s *ReviewService) UpdateReview(ctx context.Context, userID, productID uuid.UUID, title, reviewText string, rating int, anonymous bool) (models.ReviewDisplay, error) {
	logger := s.logger.With(
		zap.String("method", "UpdateReview"),
		zap.String("userID", userID.String()),
		zap.String("productID", productID.String()),
	)

	reviewed, err := s.db.HasUserReviewedProduct(ctx, database.HasUserReviewedProductParams{
		UserID:    userID,
		ProductID: productID,
	})
	if err != nil {
		s.logger.Error("failed to check if user already reviewed product", zap.Error(err))
		return models.ReviewDisplay{}, fmt.Errorf("failed to check existing review: %w", apperrors.ErrInternal)
	}

	if !reviewed {
		logger.Info("user attempted to update a non existent review")
		return models.ReviewDisplay{}, fmt.Errorf("user has not reviewed the product: %w", apperrors.ErrNotFound)
	}

	_, err = s.db.UpdateUserReview(ctx, database.UpdateUserReviewParams{
		Title: sql.NullString{
			String: title,
			Valid:  true,
		},
		UpdatedAt: time.Now(),
		ReviewText: sql.NullString{
			String: reviewText,
			Valid:  true,
		},
		Rating:    int32(rating),
		UserID:    userID,
		ProductID: productID,
		Anonymous: anonymous,
	})
	if err != nil {
		logger.Error("failed to update review", zap.Error(err))
		return models.ReviewDisplay{}, fmt.Errorf("failed to update review: %w", apperrors.ErrInternal)
	}

	updatedReview, err := s.GetReviewByUserAndProduct(ctx, userID, productID)
	if err != nil {
		logger.Error("failed to retrieve updated review", zap.Error(err))
		return models.ReviewDisplay{}, fmt.Errorf("failed to retrieve updated review: %w", apperrors.ErrInternal)
	}

	logger.Info("review updated successfully")
	return updatedReview, nil
}
