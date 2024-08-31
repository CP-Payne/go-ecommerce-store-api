package handlers

import (
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

type UserHandler struct {
	srv    *service.UserService
	logger *zap.Logger
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	logger := config.GetLogger()
	return &UserHandler{
		srv:    srv,
		logger: logger,
	}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	formParams := struct {
		email           string
		name            string
		password        string
		confirmPassword string
	}{
		email:           r.FormValue("email"),
		name:            r.FormValue("name"),
		password:        r.FormValue("password"),
		confirmPassword: r.FormValue("confirm_password"),
	}

	if formParams.password != formParams.confirmPassword {
		utils.RespondWithError(w, http.StatusBadRequest, "passwords do not match")
		return
	}
	// TODO: Perform validation on email, name and password
	//

	user, err := h.srv.CreateUser(r.Context(), formParams.email, formParams.name, formParams.password)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			utils.RespondWithError(w, http.StatusConflict, "Email already exist")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to register user")
		h.logger.Info("failed to register user", zap.Error(err), zap.String("email", formParams.email))
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

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	formParams := struct {
		email    string
		password string
	}{
		email:    r.FormValue("email"),
		password: r.FormValue("password"),
	}

	user, err := h.srv.GetUserByEmail(r.Context(), formParams.email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid credentials")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "failed login")
		return
	}
	err = hashing.CheckPasswordHash(formParams.password, user.HashedPassword)
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

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}
