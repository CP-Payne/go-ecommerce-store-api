package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/CP-Payne/ecomstore/internal/utils"
	"github.com/CP-Payne/ecomstore/internal/utils/apperrors"
	"github.com/CP-Payne/ecomstore/internal/utils/hashing"
	"go.uber.org/zap"
)

type AuthHandler struct {
	srv    *service.UserService
	logger *zap.Logger
}

func NewAuthHandler(srv *service.UserService) *AuthHandler {
	logger := config.GetLogger()
	return &AuthHandler{
		srv:    srv,
		logger: logger,
	}
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	type inputParams struct {
		Email           string `json:"email"`
		Name            string `json:"name"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	params := &inputParams{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	if params.Password != params.ConfirmPassword {
		utils.RespondWithError(w, http.StatusBadRequest, "passwords do not match")
		return
	}
	// TODO: Perform validation on email, name and password

	user, err := h.srv.CreateUser(r.Context(), params.Email, params.Name, params.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			utils.RespondWithError(w, http.StatusConflict, "Email already exist")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to register user")
		h.logger.Info("failed to register user", zap.Error(err), zap.String("email", params.Email))
		return
	}

	token := config.MakeToken(user.Email, user.ID)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Value:    token,
	})

	response := struct {
		Msg string `json:"msg"`
	}{
		Msg: "Registration successfull",
	}

	utils.RespondWithJson(w, http.StatusCreated, &response)
	// http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	type inputParams struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := &inputParams{}

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	user, err := h.srv.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid credentials")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "failed login")
		return
	}
	err = hashing.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid credentials")
		return
	}

	token := config.MakeToken(user.Email, user.ID)

	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Value:    token,
	})
	response := struct {
		Msg string `json:"msg"`
	}{
		Msg: "Login successfull",
	}

	utils.RespondWithJson(w, http.StatusOK, &response)
	// http.Redirect(w, r, "/home", http.StatusSeeOther)
}
