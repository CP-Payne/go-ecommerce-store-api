package models

import (
	"time"

	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/google/uuid"
)

type User struct {
	CreatedAt      time.Time `json:"createdAt,omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	Name           string    `json:"name"`
	ID             uuid.UUID `json:"id"`
}

// Database User to User mappings
func DatabaseUserToUser(user database.User) User {
	return User{
		ID:             user.ID,
		Email:          user.Email,
		Name:           NullStringToString(user.Name),
		HashedPassword: user.HashedPassword,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}
