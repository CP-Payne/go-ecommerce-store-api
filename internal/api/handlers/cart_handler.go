package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/CP-Payne/ecomstore/internal/utils"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CartHandler struct {
	srv    *service.CartService
	logger *zap.Logger
}

func NewCartHandler(srv *service.CartService) *CartHandler {
	logger := config.GetLogger()
	return &CartHandler{
		srv:    srv,
		logger: logger,
	}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "GetCart"))

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

	cart, err := h.srv.GetCart(ctx, userID)
	if err != nil {
		logger.Error("failed to retrieve user cart", zap.Error(err), zap.String("userID", userID.String()))
		utils.RespondWithError(w, http.StatusNotFound, "Failed to retrieve user cart")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, cart)
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "AddToCart"))

	type CartInput struct {
		ProductID uuid.UUID `json:"productId"`
		Quantity  int       `json:"quantity"`
	}

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

	var cartInput CartInput

	if err := json.NewDecoder(r.Body).Decode(&cartInput); err != nil {
		logger.Warn("failed to decode request body", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if cartInput.Quantity < 1 {
		logger.Warn("user provided invalid product quantity", zap.String("userID", userID.String()), zap.Int("quantity", cartInput.Quantity))
		utils.RespondWithError(w, http.StatusBadRequest, "Quantity must be greater than 0")
		return
	}

	err = h.srv.AddToCart(ctx, userID, cartInput.ProductID, cartInput.Quantity)
	if err != nil {
		logger.Error("failed to add item to cart", zap.Error(err), zap.String("userID", userID.String()), zap.String("productID", cartInput.ProductID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add item to cart")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"message": "Succesfully added item to cart",
	})
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "RemoveFromCart"))

	type CartInput struct {
		ProductID uuid.UUID `json:"productId"`
	}

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

	var cartInput CartInput
	if err := json.NewDecoder(r.Body).Decode(&cartInput); err != nil {
		logger.Warn("failed to decode request body", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.srv.RemoveFromCart(ctx, userID, cartInput.ProductID)
	if err != nil {
		logger.Error("failed to remove item from cart", zap.Error(err), zap.String("userID", userID.String()), zap.String("productID", cartInput.ProductID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to remove item from cart")
		return
	}
	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"message": "Item removed succesfully",
	})
}

func (h *CartHandler) ReduceFromCart(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "ReduceFromCart"))

	type CartInput struct {
		ProductID uuid.UUID `json:"productId"`
	}

	var cartInput CartInput

	if err := json.NewDecoder(r.Body).Decode(&cartInput); err != nil {
		logger.Warn("failed to decode request body", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

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

	err = h.srv.ReduceFromCart(r.Context(), userID, cartInput.ProductID, 1)
	if err != nil {
		logger.Error("failed to reduce cart item quantity", zap.Error(err), zap.String("userID", userID.String()), zap.String("productID", cartInput.ProductID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to reduce cart item quantity")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"message": "Succesfully reduced cart item quantity",
	})
}
