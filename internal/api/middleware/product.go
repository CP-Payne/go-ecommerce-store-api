package middleware

import (
	"context"
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/CP-Payne/ecomstore/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type key string

const ProductIDKey key = "productID"

// ProductMiddleware is a middleware that validates the product ID in the URL and check if the product exists
func ProductMiddleware(srvProduct *service.ProductService, logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			logger := logger.With(zap.String("middleware", "ProductMiddleware"))
			// Extract product ID from URL
			strProductID := chi.URLParam(r, "id")
			// Parse the product ID
			productID, err := uuid.Parse(strProductID)
			if err != nil {
				logger.Warn("invalid product id", zap.Error(err), zap.String("productID", strProductID))
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
				return
			}

			// Check if the product exists
			productExists, err := srvProduct.ProductExists(r.Context(), productID)
			if err != nil {
				logger.Error("failed to check product existence", zap.Error(err), zap.String("productID", productID.String()))
				utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process request")
				return
			}

			if !productExists {
				logger.Warn("product does not exist", zap.String("productID", productID.String()))
				utils.RespondWithError(w, http.StatusNotFound, "Product not found")
				return
			}

			// Store the productID in the request context
			ctx := context.WithValue(r.Context(), ProductIDKey, productID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
