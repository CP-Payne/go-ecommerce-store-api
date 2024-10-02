package handlers

import (
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/CP-Payne/ecomstore/internal/utils"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type OrderHandler struct {
	srv    *service.OrderService
	logger *zap.Logger
}

func NewOrderHandler(srv *service.OrderService) *OrderHandler {
	logger := config.GetLogger()
	return &OrderHandler{
		srv:    srv,
		logger: logger,
	}
}

func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "GetUserOrders"))

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

	userOrders, err := h.srv.GetUserOrders(ctx, userID)
	if err != nil {
		logger.Error("failed to retrieve user orders", zap.Error(err), zap.String("userID", userID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user orders")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, userOrders)
}
