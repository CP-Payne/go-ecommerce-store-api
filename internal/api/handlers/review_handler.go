package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/CP-Payne/ecomstore/internal/config"
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
	Anonymous  bool   `json:"anonymous"`
}

func (h *ReviewHandler) GetProductReviews(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "GetProductReviews"))

	strProductID := chi.URLParam(r, "id")
	productID, err := uuid.Parse(strProductID)
	if err != nil {
		logger.Warn("invalid product id", zap.Error(err), zap.String("productID", strProductID))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	productExists, err := h.srvProduct.ProductExists(ctx, productID)
	if err != nil {
		logger.Error("failed to check product existence", zap.Error(err), zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	if !productExists {
		logger.Warn("product does not exist", zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	productReviews, err := h.srvReview.GetProductReviews(ctx, productID)
	if err != nil {
		logger.Error("failed to retrieve product reviews", zap.Error(err), zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve product reviews")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, productReviews)
}

func (h *ReviewHandler) AddReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "AddReview"))

	strProductID := chi.URLParam(r, "id")
	productID, err := uuid.Parse(strProductID)
	if err != nil {
		logger.Warn("invalid product id", zap.Error(err), zap.String("productID", strProductID))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	productExists, err := h.srvProduct.ProductExists(ctx, productID)
	if err != nil {
		logger.Error("failed to check product existence", zap.Error(err), zap.String("product ID", productID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed process request")
		return
	}

	if !productExists {
		logger.Warn("product does not exist", zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

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

	var reviewInput ReviewInput
	if err := json.NewDecoder(r.Body).Decode(&reviewInput); err != nil {
		logger.Warn("failed to decode request body", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validateReviewInput(reviewInput); err != nil {
		logger.Warn("invalid review input", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.srvReview.PostReview(ctx, reviewInput.Title, reviewInput.ReviewText, reviewInput.Rating, reviewInput.Anonymous, productID, userID)
	if err != nil {
		if errors.Is(err, apperrors.ErrConflict) {
			logger.Info("user has already reviewed the product", zap.String("userID", userID.String()), zap.String("productID", productID.String()))
			utils.RespondWithError(w, http.StatusConflict, "User has already reviewed the product")
			return
		}
		logger.Error("failed to post review", zap.Error(err), zap.String("userID", userID.String()), zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to add review")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, map[string]interface{}{
		"message": "Review added successfully",
	})
}

func validateReviewInput(input ReviewInput) error {
	if len(input.Title) > 30 {
		return errors.New("title must be less than 30 characters")
	}
	if input.Rating < 1 || input.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}
	return nil
}

func (h *ReviewHandler) GetUserReviewForProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "GetUserReviewForProduct"))

	// TODO: Why not turn below into middleware?
	strProductID := chi.URLParam(r, "productID")
	productID, err := uuid.Parse(strProductID)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product id")
		return
	}

	productExists, err := h.srvProduct.ProductExists(ctx, productID)
	if err != nil {
		logger.Error("failed to check product existence", zap.Error(err), zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	if !productExists {
		logger.Warn("product does not exist", zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

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

	review, err := h.srvReview.GetReviewByUserAndProduct(r.Context(), userID, productID)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			logger.Warn("user review does not exist", zap.String("userID", userID.String()), zap.String("productID", productID.String()))
			utils.RespondWithError(w, http.StatusNotFound, "User review not found")
			return
		}
		logger.Error("failed to retrieve user review", zap.Error(err), zap.String("userID", userID.String()), zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user review")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, review)
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "DeleteReview"))

	strProductID := chi.URLParam(r, "id")
	productID, err := uuid.Parse(strProductID)
	if err != nil {
		logger.Warn("invalid product id", zap.Error(err), zap.String("productID", strProductID))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	productExists, err := h.srvProduct.ProductExists(ctx, productID)
	if err != nil {
		logger.Error("failed to check product existence", zap.Error(err), zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	if !productExists {
		logger.Warn("product does not exist", zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

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

	err = h.srvReview.DeleteReview(ctx, userID, productID)
	if err != nil {
		logger.Error("failed to delete review", zap.Error(err), zap.String("userID", userID.String()), zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to delete review")
		return
	}

	utils.RespondWithJson(w, http.StatusOK, map[string]interface{}{
		"message": "Review deleted successfully",
	})
}

func (h *ReviewHandler) UpdateUserReview(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger := h.logger.With(zap.String("handler", "UpdateUserReview"))

	strProductID := chi.URLParam(r, "id")
	productID, err := uuid.Parse(strProductID)
	if err != nil {
		logger.Warn("invalid product id", zap.Error(err), zap.String("productID", strProductID))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	productExists, err := h.srvProduct.ProductExists(ctx, productID)
	if err != nil {
		logger.Error("failed to check product existence", zap.Error(err), zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process request")
		return
	}

	if !productExists {
		logger.Warn("product does not exist", zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

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

	var reviewInput ReviewInput
	reviewInput.Anonymous = true

	if err := json.NewDecoder(r.Body).Decode(&reviewInput); err != nil {
		logger.Warn("failed to decode request body", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := validateReviewInput(reviewInput); err != nil {
		logger.Warn("invalid review input", zap.Error(err))
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	review, err := h.srvReview.UpdateReview(ctx, userID, productID, reviewInput.Title, reviewInput.ReviewText, reviewInput.Rating, reviewInput.Anonymous)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			logger.Info("review not found", zap.String("userID", userID.String()), zap.String("productID", productID.String()))
			utils.RespondWithError(w, http.StatusNotFound, "User has not reviewed the product")
			return
		}
		logger.Error("failed to update review", zap.Error(err), zap.String("userID", userID.String()), zap.String("productID", productID.String()))
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to update review")
		return
	}

	utils.RespondWithJson(w, http.StatusCreated, map[string]interface{}{
		"message": "Review updated successfully",
		"review":  review,
	})
}
