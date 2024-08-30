package service

import (
	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
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

// TODO: Implement domain models. Implement some type of error handling package
func (s *UserService) GetUserByEmail(email string) {
}

func (s *UserService) CreateUser(email, hashedPassword string)
