package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/domain/user"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/CP-Payne/ecomstore/internal/utils"
	"github.com/CP-Payne/ecomstore/internal/utils/apperrors"
	"github.com/CP-Payne/ecomstore/internal/utils/hashing"
	"github.com/CP-Payne/ecomstore/pkg/errsx"
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

type RegistrationInput struct {
	Email           string `json:"email"`
	Name            string `json:"name"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "RegisterUser"))

	var input RegistrationInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Warn("failed to decode request body", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if input.Password != input.ConfirmPassword {
		utils.RespondWithError(w, http.StatusBadRequest, "Passwords do not match")
		return
	}

	errs := input.validateRegistrationInput()

	if errs != nil {
		// If the map is not nil, we have some errors.
		utils.RespondWithJson(w, http.StatusBadRequest, errs)
		return
	}

	user, err := h.srv.CreateUser(ctx, input.Email, input.Name, input.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			logger.Info("user email already exists", zap.String("email", input.Email))
			utils.RespondWithError(w, http.StatusConflict, "Email already exists")
			return
		}
		h.logger.Info("failed to register user", zap.Error(err), zap.String("email", input.Email))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to register user")
		return
	}

	token := config.MakeToken(user.Email, user.ID)
	setCookie(w, token)

	utils.RespondWithJson(w, http.StatusCreated, map[string]interface{}{
		"message": "Registration successfull",
		"email":   user.Email,
		"userId":  user.ID.String(),
	})
}

func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "LoginUser"))

	var params LoginInput

	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		logger.Warn("failed to decode request body", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.srv.GetUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			logger.Info("login failed", zap.Error(err))
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid credentials")
			return
		}
		logger.Error("failed to log in user", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to authenticate user")
		return
	}
	err = hashing.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		logger.Info("login failed", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid credentials")
		return
	}

	token := config.MakeToken(user.Email, user.ID)
	setCookie(w, token)

	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"message": "Login successfull",
		"userId":  user.ID.String(),
	})
}

func setCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		HttpOnly: true,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		SameSite: http.SameSiteLaxMode,
		Name:     "jwt",
		Value:    token,
	})
}

func (ri *RegistrationInput) validateRegistrationInput() errsx.Map {
	var err error
	var errs errsx.Map

	_, err = user.ValidateEmail(ri.Email)
	if err != nil {
		errs.Set("email", err)
	}

	_, err = user.ValidateName(ri.Name)
	if err != nil {
		errs.Set("name", err)
	}
	_, err = user.ValidatePassword(ri.Password)
	if err != nil {
		errs.Set("password", err)
	}

	return errs
}
