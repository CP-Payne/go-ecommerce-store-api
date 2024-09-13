package api

import (
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/api/handlers"
	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"

	cmid "github.com/CP-Payne/ecomstore/internal/api/middleware"
)

func SetupRouter(cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Use(cmid.CorsMiddleware)

	userSrv := service.NewUserService(cfg.DB)
	productSrv := service.NewProductService(cfg.DB)
	reviewSrv := service.NewReviewService(cfg.DB)
	cartSrv := service.NewCartService(cfg.DB)

	authHandler := handlers.NewAuthHandler(userSrv)
	productHandler := handlers.NewProductHandler(productSrv)
	reviewHander := handlers.NewReviewHandler(reviewSrv, productSrv)
	cartHandler := handlers.NewCartHandler(cartSrv)
	// TODO: if logged in and logging request is sent, redirect user to home page or profile

	r.Group(func(r chi.Router) {
		r.Post("/register", authHandler.RegisterUser)
		r.Post("/login", authHandler.LoginUser)

		r.Get("/products", productHandler.GetAllProducts)
		r.Get("/products/{id}", productHandler.GetProduct)

		r.Get("/products/categories", productHandler.GetProductCategories)
		r.Get("/products/categories/{id}", productHandler.GetProductsByCategory)

		r.Get("/products/{id}/reviews", reviewHander.GetProductReviews)
		r.Get("/products/{productID}/reviews/{userID}", reviewHander.GetUserReviewForProduct)
	})

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.GetTokenAuth()))
		r.Use(jwtauth.Authenticator)
		r.Post("/products/{id}/reviews", reviewHander.AddReview)
		r.Patch("/products/{id}/reviews", reviewHander.UpdateUserReview)
		r.Delete("/products/{id}/reviews", reviewHander.DeleteReview)

		r.Get("/cart", cartHandler.GetCart)
		r.Post("/cart/add", cartHandler.AddToCart)
		r.Post("/cart/remove", cartHandler.RemoveFromCart)
	})

	r.Get("/home", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This will be the home page"))
	}))

	return r
}
