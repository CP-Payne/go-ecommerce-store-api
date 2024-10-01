package handlers

import (
	"errors"
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/CP-Payne/ecomstore/internal/utils"
	"github.com/CP-Payne/ecomstore/internal/utils/apperrors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProductHandler struct {
	srv    *service.ProductService
	logger *zap.Logger
}

func NewProductHandler(srv *service.ProductService) *ProductHandler {
	logger := config.GetLogger()
	return &ProductHandler{
		srv:    srv,
		logger: logger,
	}
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "GetProduct"))

	strID := chi.URLParam(r, "id")
	id, err := uuid.Parse(strID)
	if err != nil {
		logger.Warn("invalid product id", zap.Error(err), zap.String("productID", strID))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.srv.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			logger.Info("product not found", zap.Error(err), zap.String("productID", strID))
			utils.RespondWithError(w, http.StatusNotFound, "Product not found")
			return
		}

		logger.Info("failed to retrieve product", zap.Error(err), zap.String("productID", strID))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve product")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, product)
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "GetAllProducts"))

	products, err := h.srv.GetAllProducts(ctx)
	if err != nil {
		logger.Error("failed to retrieve product list", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve products")
		return
	}
	utils.RespondWithJson(w, http.StatusOK, products)
}

func (h *ProductHandler) GetProductCategories(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "GetProductCategories"))

	categories, err := h.srv.GetProductCategories(ctx)
	if err != nil {
		logger.Error("failed to retrieve product categories", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve product categories")
		return
	}
	utils.RespondWithJson(w, http.StatusOK, categories)
}

func (h *ProductHandler) GetProductsByCategory(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "GetProductsByCategory"))

	strID := chi.URLParam(r, "id")
	id, err := uuid.Parse(strID)
	if err != nil {

		logger.Warn("invalid category id", zap.Error(err), zap.String("productID", strID))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	products, err := h.srv.GetProductsByCategory(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			logger.Info("no products for the provided category found", zap.Error(err))
			utils.RespondWithError(w, http.StatusNotFound, "No products for provided category found")
			return
		}
		logger.Error("failed to retrieve products for category", zap.Error(err), zap.String("categoryID", strID))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve products for category")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, products)
}
