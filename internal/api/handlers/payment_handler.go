package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/CP-Payne/ecomstore/internal/utils"
	"github.com/CP-Payne/ecomstore/internal/utils/apperrors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	srvProduct *service.ProductService
	srvPayment *service.PaymentService
	logger     *zap.Logger
}

func NewPaymentHandler(srvProduct *service.ProductService, srvPayment *service.PaymentService) *PaymentHandler {
	logger := config.GetLogger()
	return &PaymentHandler{
		srvProduct: srvProduct,
		srvPayment: srvPayment,
		logger:     logger,
	}
}

func (h *PaymentHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	type inputParams struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	}

	params := &inputParams{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	id, err := uuid.Parse(params.ProductID)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	if params.Quantity < 1 {
		params.Quantity = 1
	}

	product, err := h.srvProduct.GetProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "product not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to retrieve product")
		h.logger.Info("failed to retrieve product", zap.Error(err), zap.String("productID", params.ProductID))
		return
	}

	if params.Quantity > product.Stock {
		utils.RespondWithError(w, http.StatusBadRequest, "not enough stock")
		return
	}

	// TODO: Get user ID and save order to database

	// TODO: Input validation
	orderResult, err := h.srvPayment.CreateOrder(r.Context(), &product, params.Quantity)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, orderResult)
}

func (h *PaymentHandler) CaptureOrder(w http.ResponseWriter, r *http.Request) {
	orderID := r.URL.Query().Get("token")
	if orderID == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "no token (orderID) query found")
		return
	}

	err := h.srvPayment.CaptureOrder(r.Context(), orderID)
	if err != nil {
		// TODO: Check what the error is and return the appropriate message
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to complete payment")
		return
	}

	type successResponse struct {
		Msg string `json:"msg"`
	}

	utils.RespondWithJson(w, http.StatusOK, successResponse{
		Msg: "purchase successfull",
	})
}
