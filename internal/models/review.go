package models

import (
	"time"

	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/google/uuid"
)

type Review struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	ReviewText string    `json:"reviewText,omitempty"`
	Rating     int       `json:"rating,omitempty"`
	ProductID  uuid.UUID `json:"productId,omitempty"`
	UserID     uuid.UUID `json:"userId,omitempty"`
	Deleted    bool      `json:"deleted,omitempty"`
	CreatedAt  time.Time `json:"createdAt,omitempty"`
	UpdatedAt  time.Time `json:"updatedAt,omitempty"`
}

// Database Review to Review mappings
func DatabaseReviewToReview(dbReview database.Review) Review {
	return Review{
		ID:         dbReview.ID,
		Title:      NullStringToString(dbReview.Title),
		ReviewText: NullStringToString(dbReview.ReviewText),
		Rating:     int(dbReview.Rating),
		ProductID:  dbReview.ProductID,
		UserID:     dbReview.UserID,
		Deleted:    dbReview.Deleted,
		CreatedAt:  dbReview.CreatedAt,
		UpdatedAt:  dbReview.UpdatedAt,
	}
}

func DatabaseReviewsToReviews(dbReviews []database.Review) []Review {
	reviews := make([]Review, 0, len(dbReviews))
	for _, dbRev := range dbReviews {
		reviews = append(reviews, DatabaseReviewToReview(dbRev))
	}
	return reviews
}
