package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/config"
	"github.com/CP-Payne/ecomstore/internal/models"
	"github.com/CP-Payne/ecomstore/internal/service"
	"github.com/CP-Payne/ecomstore/internal/utils"
	"github.com/CP-Payne/ecomstore/internal/utils/apperrors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ReviewHandler struct {
	srv    *service.ReviewService
	logger *zap.Logger
}

func NewReviewHandler(srv *service.ReviewService) *ReviewHandler {
	logger := config.GetLogger()
	return &ReviewHandler{
		srv:    srv,
		logger: logger,
	}
}

type ReviewInput struct {
	Title      string `json:"title"`
	ReviewText string `json:"reviewText"`
	Rating     int    `json:"rating"`
}

func (h *ReviewHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	// Get productID from url parameter
	strProductID := chi.URLParam(r, "id")
	productID, err := uuid.Parse(strProductID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid product id")
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
	}

	reviewInput := &ReviewInput{}
	if err := json.NewDecoder(r.Body).Decode(reviewInput); err != nil {
		h.logger.Error("failed to decode json body", zap.Error(err))
		utils.RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	if len(reviewInput.Title) > 30 {
		utils.RespondWithError(w, http.StatusBadRequest, "title size must be less than 30 characters")
		return
	}

	if reviewInput.Rating > 5 {
		reviewInput.Rating = 5
	} else if reviewInput.Rating < 1 {
		reviewInput.Rating = 1
	}

	review, err := h.srv.PostReview(r.Context(), reviewInput.Title, reviewInput.ReviewText, reviewInput.Rating, productID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			utils.RespondWithError(w, http.StatusConflict, "user already reviewed the product")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "could not add review")
		return
	}

	type successResponse struct {
		Msg    string        `json:"msg"`
		Review models.Review `json:"review"`
		// TODO: return user's name instead of the ID
		// TODO: Do not return entire review object
	}

	utils.RespondWithJson(w, http.StatusCreated, successResponse{
		Msg:    "review added successfully",
		Review: review,
	})
}
