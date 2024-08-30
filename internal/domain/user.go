package domain

import "time"

type User struct {
	CreatedAt      time.Time `json:"created_at,omitempty"`
	UpdatedAt      time.Time `json:"updated_at,omitempty"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	ID             int       `json:"id"`
}
