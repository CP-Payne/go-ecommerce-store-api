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

type UserHandler struct {
	srvUser *service.UserService
	logger  *zap.Logger
}

func NewUserHandler(userSrv *service.UserService) *UserHandler {
	logger := config.GetLogger()
	return &UserHandler{
		srvUser: userSrv,
		logger:  logger,
	}
}

func (h *UserHandler) GetUserDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "GetUserDetails"))

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

	user, err := h.srvUser.GetUserProfile(ctx, userID)
	if err != nil {
		logger.Error("failed to retrieve user information", zap.Error(err), zap.String("userID", userID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user information")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, user)
}
