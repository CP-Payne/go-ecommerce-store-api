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
	srvOrder   *service.OrderService
	logger     *zap.Logger
}

func NewPaymentHandler(srvProduct *service.ProductService, srvPayment *service.PaymentService, srvCart *service.CartService, srvOrder *service.OrderService) *PaymentHandler {
	logger := config.GetLogger()
	return &PaymentHandler{
		srvProduct: srvProduct,
		srvPayment: srvPayment,
		srvCart:    srvCart,
		srvOrder:   srvOrder,
		logger:     logger,
	}
}

func (h *PaymentHandler) CreateOrderCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "CreateOrderCart"))

	_, claims, _ := jwtauth.FromContext(ctx)
	strUserID, ok := claims["id"].(string)
	if !ok {
		logger.Error("user id not found in token claims")
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid authentication")
		return
	}
	userID, err := uuid.Parse(strUserID)
	if err != nil {
		logger.Error("failed to parse user id", zap.Error(err), zap.String("userID", strUserID))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	// cart, err := h.srvCart.GetCartByID(r.Context(), userID, cartID)
	cart, err := h.srvCart.GetCart(ctx, userID)
	if err != nil {
		logger.Info("failed to retrieve user cart", zap.Error(err), zap.String("UserID", userID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}
	if len(cart.Items) <= 0 {
		logger.Info("user attempted checkout an empty cart", zap.String("userID", userID.String()))
		utils.RespondWithError(w, http.StatusBadRequest, "cart is empty")
		return
	}

	// Check product quantity and cart quantity
	for _, ci := range cart.Items {
		product, err := h.srvProduct.GetProduct(ctx, ci.ProductID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process request")
			logger.Info("failed to retrieve product during checkout", zap.Error(err), zap.String("productID", ci.ProductID.String()))
			return
		}

		if ci.Quantity > product.Stock {
			logger.Warn("insufficient stock to create order", zap.String("ProductID", ci.ProductID.String()), zap.String("CartID", cart.ID.String()), zap.String("UserID", userID.String()))
			utils.RespondWithError(w, http.StatusBadRequest, "Not enough stock")
			return
		}
	}

	order, err := h.srvOrder.CreateOrder(ctx, cart, false)
	if err != nil {
		logger.Error("failed to create order for user", zap.Error(err), zap.String("userID", userID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	orderResult, err := h.srvPayment.CreateProcessorOrder(ctx, &order)
	if err != nil {
		logger.Error("failed to create processor order for user", zap.Error(err), zap.String("orderID", order.ID.String()), zap.String("userID", userID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	logger.Info("succesfully created order for user", zap.String("userID", userID.String()))
	utils.RespondWithJson(w, http.StatusOK, orderResult)
}

func (h *PaymentHandler) CreateOrderProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "CreateOrderProduct"))

	_, claims, _ := jwtauth.FromContext(ctx)
	strUserID, ok := claims["id"].(string)
	if !ok {
		logger.Error("user id not found in token claims")
		utils.RespondWithError(w, http.StatusUnauthorized, "Invalid authentication")
		return
	}
	userID, err := uuid.Parse(strUserID)
	if err != nil {
		logger.Error("failed to parse user id", zap.Error(err), zap.String("userID", strUserID))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	// cart, err := h.srvCart.GetCartByID(r.Context(), userID, cartID)

	type inputParams struct {
		ProductID string `json:"productId"`
		Quantity  int    `json:"quantity"`
	}

	params := &inputParams{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Warn("failed to decode request body", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	id, err := uuid.Parse(params.ProductID)
	if err != nil {
		logger.Warn("invalid product id", zap.Error(err), zap.String("productID", params.ProductID))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if params.Quantity < 1 {
		logger.Info("user provided invalid order quantity", zap.Error(err), zap.Int("quantity", params.Quantity))
		utils.RespondWithError(w, http.StatusBadRequest, "Quantity cannot be less than 1")
		return
	}

	product, err := h.srvProduct.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			logger.Warn("user attempted to purchase a product that does not exist")
			utils.RespondWithError(w, http.StatusNotFound, "Product not found")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process request")
		logger.Info("failed to retrieve product during checkout", zap.Error(err), zap.String("productID", params.ProductID))
		return
	}

	if params.Quantity > product.Stock {
		logger.Warn("insufficient stock to create order")
		utils.RespondWithError(w, http.StatusBadRequest, "Not enough stock")
		return
	}

	// Create temporary cart
	tempCart := h.srvCart.CreateTemporaryProductCart(ctx, userID, product, params.Quantity)

	order, err := h.srvOrder.CreateOrder(ctx, tempCart, true)
	if err != nil {
		logger.Warn("failed to create order", zap.Error(err), zap.String("userID", userID.String()), zap.String("cartID", tempCart.ID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	orderResult, err := h.srvPayment.CreateProcessorOrder(r.Context(), &order)
	if err != nil {
		logger.Warn("failed to create processor order", zap.Error(err), zap.String("userID", userID.String()), zap.String("orderID", order.ID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	logger.Info("order successfully created")
	utils.RespondWithJson(w, http.StatusOK, orderResult)
}

func (h *PaymentHandler) CaptureOrder(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "CaptureOrder"))

	orderID := r.URL.Query().Get("token")
	if orderID == "" {
		logger.Warn("user did not provide a token to complete purchase")
		utils.RespondWithError(w, http.StatusBadRequest, "Token not provided")
		return
	}

	err := h.srvPayment.CaptureOrder(ctx, orderID)
	if err != nil {
		logger.Warn("failed to capture order", zap.Error(err), zap.String("orderID", orderID))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to complete payment")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"message": "Purchase succesfull",
	})
}
