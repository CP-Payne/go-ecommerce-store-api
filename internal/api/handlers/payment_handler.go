package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/CP-Payne/ecomstore/internal/utils"
	"github.com/CP-Payne/ecomstore/internal/utils/apperrors"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PaymentHandler struct {
	srvProduct *service.ProductService
	srvPayment *service.PaymentService
	srvCart    *service.CartService
	logger     *zap.Logger
}

func NewPaymentHandler(srvProduct *service.ProductService, srvPayment *service.PaymentService, srvCart *service.CartService) *PaymentHandler {
	logger := config.GetLogger()
	return &PaymentHandler{
		srvProduct: srvProduct,
		srvPayment: srvPayment,
		srvCart:    srvCart,
		logger:     logger,
	}
}

func (h *PaymentHandler) CreateProductOrder(w http.ResponseWriter, r *http.Request) {
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
	orderResult, err := h.srvPayment.CreateProductOrder(r.Context(), &product, params.Quantity)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, orderResult)
}

func (h *PaymentHandler) CreateCartOrder(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	strUserID := ""
	if claims["id"] != nil {
		strUserID = claims["id"].(string)
	} else {
		utils.RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	userID, err := uuid.Parse(strUserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// cart, err := h.srvCart.GetCartByID(r.Context(), userID, cartID)
	cart, err := h.srvCart.GetCart(r.Context(), userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "cart not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to retrieve cart")
		h.logger.Info("failed to retrieve user cart", zap.Error(err), zap.String("UserID", userID.String()))
		return
	}
	if len(cart.Items) <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "cart is empty")
		return
	}

	// TODO: Input validation
	orderResult, err := h.srvPayment.CreateCartOrder(r.Context(), &cart)
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
