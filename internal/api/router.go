package api

import (
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/api/handlers"
	cmid "github.com/CP-Payne/ecomstore/internal/api/middleware"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(db *database.Queries) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// r.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"http://localhost:3000"},
	// 	AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "OPTION"},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: true,
	// 	MaxAge:           300, // Maximum value not ignored by any of major browsers
	// }))

	r.Use(cmid.CorsMiddleware)

	userSrv := service.NewUserService(db)
	productSrv := service.NewProductService(db)

	authHandler := handlers.NewAuthHandler(userSrv)
	productHandler := handlers.NewProductHandler(productSrv)
	// TODO: if logged in and logging request is sent, redirect user to home page or profile

	r.Group(func(r chi.Router) {
		r.Post("/register", authHandler.RegisterUser)
		r.Post("/login", authHandler.LoginUser)

		r.Get("/product/{id}", productHandler.GetProduct)
	})

	r.Get("/home", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This will be the home page"))
	}))

	return r
}
