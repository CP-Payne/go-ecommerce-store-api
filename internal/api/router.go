package api

import (
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/api/handlers"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func SetupRouter(db *database.Queries) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "OPTION"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	userSrv := service.NewUserService(db)
	userHandler := handlers.NewUserHandler(userSrv)

	r.Group(func(r chi.Router) {
		r.Post("/register", userHandler.RegisterUser)
		r.Post("/login", userHandler.LoginUser)
	})

	r.Get("/home", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This will be the home page"))
	}))

	return r
}
