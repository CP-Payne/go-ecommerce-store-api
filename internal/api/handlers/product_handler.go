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
	strID := chi.URLParam(r, "id")
	id, err := uuid.Parse(strID)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.srv.GetProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "product not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to retrieve product")
		h.logger.Info("failed to retrieve product", zap.Error(err), zap.String("productID", strID))
		return
	}

	utils.RespondWithJson(w, http.StatusOK, product)
}

func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.srv.GetAllProducts(r.Context())
	if err != nil {
		h.logger.Error("failed to respond with product list", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to retrieve products")
		return
	}
	utils.RespondWithJson(w, http.StatusOK, products)
}
