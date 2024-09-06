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
	srvReview  *service.ReviewService
	srvProduct *service.ProductService
	logger     *zap.Logger
}

func NewReviewHandler(srvReview *service.ReviewService, srvProduct *service.ProductService) *ReviewHandler {
	logger := config.GetLogger()
	return &ReviewHandler{
		srvReview:  srvReview,
		srvProduct: srvProduct,
		logger:     logger,
	}
}

type ReviewInput struct {
	Title      string `json:"title"`
	ReviewText string `json:"reviewText"`
	Rating     int    `json:"rating"`
}

func (h *ReviewHandler) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	strProductID := chi.URLParam(r, "id")
	productID, err := uuid.Parse(strProductID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid product id")
		// TODO: Test remove below
		return
	}

	productExists, err := h.srvProduct.ProductExists(r.Context(), productID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to retrieve product reviews")
		return
	}

	if !productExists {

		utils.RespondWithError(w, http.StatusBadRequest, "productID provided does not exist")
		return
	}

	productReviews, err := h.srvReview.GetProductReviews(r.Context(), productID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to retrieve product reviews")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, productReviews)
}

func (h *ReviewHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	// Get productID from url parameter
	// TODO: Make sure only a user who purchased the product can review it
	strProductID := chi.URLParam(r, "id")
	productID, err := uuid.Parse(strProductID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	productExists, err := h.srvProduct.ProductExists(r.Context(), productID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not determine if product exists")
		return
	}

	if !productExists {
		utils.RespondWithError(w, http.StatusBadRequest, "productID provided does not exist")
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

	review, err := h.srvReview.PostReview(r.Context(), reviewInput.Title, reviewInput.ReviewText, reviewInput.Rating, productID, userID)
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

func (h *ReviewHandler) GetUserReviewForProduct(w http.ResponseWriter, r *http.Request) {
	// Get productID from url parameter
	strProductID := chi.URLParam(r, "productID")
	productID, err := uuid.Parse(strProductID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	productExists, err := h.srvProduct.ProductExists(r.Context(), productID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not determine if product exists")
		return
	}

	if !productExists {
		utils.RespondWithError(w, http.StatusBadRequest, "productID provided does not exist")
		return
	}

	strUserID := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(strUserID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid user id")
		return
	}
	review, err := h.srvReview.GetReviewByUserAndProduct(r.Context(), userID, productID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			utils.RespondWithError(w, http.StatusNotFound, "user review not found for product")
			return
		}
		utils.RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// TODO: return user's name instead of the ID
	// TODO: Do not return entire review object

	utils.RespondWithJson(w, http.StatusOK, review)
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	// Get productID from url parameter
	strProductID := chi.URLParam(r, "id")
	productID, err := uuid.Parse(strProductID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	productExists, err := h.srvProduct.ProductExists(r.Context(), productID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "could not determine if product exists")
		return
	}

	if !productExists {
		utils.RespondWithError(w, http.StatusBadRequest, "productID provided does not exist")
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

	err = h.srvReview.DeleteReview(r.Context(), userID, productID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to delete review")
		return
	}
	type successResponse struct {
		Msg string `json:"msg"`
	}
	utils.RespondWithJson(w, http.StatusOK, successResponse{
		Msg: "review deleted",
	})
}
