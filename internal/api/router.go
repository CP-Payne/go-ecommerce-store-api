package api

import (
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/api/handlers"
	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"

	cmid "github.com/CP-Payne/ecomstore/internal/api/middleware"
)

func SetupRouter(db *database.Queries) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	// r.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "OPTION"},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: true,
	// 	MaxAge:           300, // Maximum value not ignored by any of major browsers
	// }))

	r.Use(cmid.CorsMiddleware)

	userSrv := service.NewUserService(db)
	productSrv := service.NewProductService(db)
	reviewSrv := service.NewReviewService(db)

	authHandler := handlers.NewAuthHandler(userSrv)
	productHandler := handlers.NewProductHandler(productSrv)
	reviewHander := handlers.NewReviewHandler(reviewSrv)
	// TODO: if logged in and logging request is sent, redirect user to home page or profile

	r.Group(func(r chi.Router) {
		r.Post("/register", authHandler.RegisterUser)
		r.Post("/login", authHandler.LoginUser)

		r.Get("/products", productHandler.GetAllProducts)
		r.Get("/products/{id}", productHandler.GetProduct)

		r.Get("/products/categories", productHandler.GetProductCategories)
		r.Get("/products/categories/{id}", productHandler.GetProductsByCategory)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.GetTokenAuth()))
		r.Use(jwtauth.Authenticator)
		r.Post("/products/{id}/reviews", reviewHander.AddReview)
	})

	r.Get("/home", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This will be the home page"))
	}))

	return r
}
