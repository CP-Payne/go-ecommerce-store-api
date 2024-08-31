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

func (s *UserService) GetUserByEmail(email string) {
}

func (s *UserService) CreateUser(ctx context.Context, email, name, password string) (models.User, error) {
	hashedPassword, err := hashing.HashPassword(password)
	if err != nil {
		s.logger.Error("failed to hash user password", zap.String("user", email), zap.Error(err))
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
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
			return models.User{}, fmt.Errorf("email already exists: %w", apperrors.ErrConflict)
		}
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return models.DatabaseUserToUser(dbUser), nil
}
