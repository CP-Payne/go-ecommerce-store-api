package handlers

import (
	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/service"
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

// func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
// 	formParams := struct {
// 		email           string
// 		name            string
// 		password        string
// 		confirmPassword string
// 	}{
// 		email:           r.FormValue("email"),
// 		name:            r.FormValue("name"),
// 		password:        r.FormValue("password"),
// 		confirmPassword: r.FormValue("confirm_password"),
// 	}
// }
