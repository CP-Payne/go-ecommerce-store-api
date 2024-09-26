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

	user, err := h.srvUser.GetUserProfile(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to retrieve user information")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, user)
}
