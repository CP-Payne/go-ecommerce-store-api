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
	"github.com/CP-Payne/ecomstore/internal/utils/hashing"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	logger *zap.Logger
	db     *database.Queries
}

func NewUserService(db *database.Queries) *UserService {
	return &UserService{
		logger: config.GetLogger(),
		db:     db,
	}
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	logger := s.logger.With(
		zap.String("method", "GetUserByEmail"),
		zap.String("email", email),
	)
	dbUser, err := s.db.GetUserByEmail(ctx, email)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			logger.Error("user information for email not found", zap.Error(err))
			return models.User{}, fmt.Errorf("user information for email not found: %w", apperrors.ErrNotFound)
		}
		logger.Error("failed to retrieve user information by email", zap.Error(err))
		return models.User{}, fmt.Errorf("failed to retrieve user: %w", err)
	}
	return models.DatabaseUserToUser(dbUser), nil
}

func (s *UserService) CreateUser(ctx context.Context, email, name, password string) (models.User, error) {

	logger := s.logger.With(
		zap.String("method", "CreateUser"),
		zap.String("email", email),
	)

	hashedPassword, err := hashing.HashPassword(password)
	if err != nil {
		logger.Error("failed to hash user password", zap.String("user", email), zap.Error(err))
		return models.User{}, fmt.Errorf("failed to hash user password: %w", err)
	}

	dbUser, err := s.db.CreateUser(ctx, database.CreateUserParams{
		ID:    uuid.New(),
		Email: email,
		Name: sql.NullString{
			String: name,
			Valid:  true,
		},
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		if apperrors.IsUniqueViolation(err) {
			logger.Info("user attempted to register with email that already exists", zap.String("user", email), zap.Error(err))
			return models.User{}, fmt.Errorf("failed to create user: %w", apperrors.ErrConflict)
		}
		logger.Error("failed to create new user", zap.String("user", email), zap.Error(err))
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Info("user successfully created", zap.String("email", dbUser.Email))
	return models.DatabaseUserToUser(dbUser), nil
}

func (s *UserService) GetUserProfile(ctx context.Context, userID uuid.UUID) (models.UserProfile, error) {

	logger := s.logger.With(
		zap.String("method", "GetUserProfile"),
		zap.String("userID", userID.String()),
	)

	userDetailsRow, err := s.db.GetUserDetails(ctx, userID)
	if err != nil {
		if apperrors.IsNoRowsError(err) {
			logger.Info("user profile not found", zap.Error(err))
			return models.UserProfile{}, fmt.Errorf("failed to retrieve user profile: %w", apperrors.ErrNotFound)
		}
		logger.Error("failed to retrieve user profile", zap.Error(err))
		return models.UserProfile{}, fmt.Errorf("failed to retrieve user profile: %w", err)
	}

	return models.UserProfile{
		ID:    userDetailsRow.ID,
		Email: userDetailsRow.Email,
		Name:  sqlNullStringToString(userDetailsRow.Name),
	}, nil
}
