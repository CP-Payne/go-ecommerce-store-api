package config

import (
	"os"
	"sync"

	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
)

var (
	so        sync.Once
	tokenAuth *jwtauth.JWTAuth
)

// GetTokenAuth returns a singleton instance of a JWTAuth
func GetTokenAuth() *jwtauth.JWTAuth {
	so.Do(func() {
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			GetLogger().Fatal("JWT_SECRET cannot be empty")
		}
		tokenAuth = jwtauth.New("HS256", []byte(jwtSecret), nil)
	})

	return tokenAuth
}

func MakeToken(email string, id uuid.UUID) string {
	_, tokenString, _ := GetTokenAuth().Encode(map[string]interface{}{"email": email, "id": id})
	return tokenString
}
