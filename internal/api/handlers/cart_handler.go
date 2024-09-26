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
	// Get UserID from jwt (request context)
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

	cart, err := h.srv.GetCart(r.Context(), userID)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "couldn't get cart")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, cart)
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	type CartInput struct {
		ProductID uuid.UUID `json:"productId"`
		Quantity  int       `json:"quantity"`
	}

	cartInput := &CartInput{}

	if err := json.NewDecoder(r.Body).Decode(&cartInput); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	if cartInput.Quantity < 0 {
		cartInput.Quantity = 0
	}

	// Get UserID from jwt (request context)
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

	err = h.srv.AddToCart(r.Context(), userID, cartInput.ProductID, cartInput.Quantity)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to add item to cart")
		return
	}
	type successResponse struct {
		Msg string `json:"msg"`
	}

	utils.RespondWithJson(w, http.StatusOK, successResponse{
		Msg: "item added successfully",
	})
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	type CartInput struct {
		ProductID uuid.UUID `json:"productId"`
	}

	cartInput := &CartInput{}

	if err := json.NewDecoder(r.Body).Decode(&cartInput); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	// Get UserID from jwt (request context)
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

	err = h.srv.RemoveFromCart(r.Context(), userID, cartInput.ProductID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to add item to cart")
		return
	}
	type successResponse struct {
		Msg string `json:"msg"`
	}

	utils.RespondWithJson(w, http.StatusOK, successResponse{
		Msg: "item removed successfully",
	})
}

func (h *CartHandler) ReduceFromCart(w http.ResponseWriter, r *http.Request) {
	type CartInput struct {
		ProductID uuid.UUID `json:"productId"`
	}

	cartInput := &CartInput{}

	if err := json.NewDecoder(r.Body).Decode(&cartInput); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid parameters")
		return
	}

	// Get UserID from jwt (request context)
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

	err = h.srv.ReduceFromCart(r.Context(), userID, cartInput.ProductID, 1)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to reduce product quantity")
		return
	}
	type successResponse struct {
		Msg string `json:"msg"`
	}
	utils.RespondWithJson(w, http.StatusOK, successResponse{
		Msg: "Quantity reduced successfully",
	})
}
