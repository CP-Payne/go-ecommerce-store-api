package models

import (
	"time"

	"github.com/google/uuid"
)

type Review struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Title      string    `json:"title,omitempty"`
	ReviewText string    `json:"review_text,omitempty"`
	Rating     int       `json:"rating,omitempty"`
	ProductID  uuid.UUID `json:"product_id,omitempty"`
	UserID     uuid.UUID `json:"user_id,omitempty"`
	Deleted    bool      `json:"deleted,omitempty"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

// Database Category to Category mappings
// func DatabaseCategoryToCategory(category database.Category) Category {
// 	return Category{
// 		ID:          category.ID,
// 		Name:        category.Name,
// 		Description: NullStringToString(category.Description),
// 	}
// }
//
// func DatabaseCategoriesToCategories(dbReviews []database.Category) []Category {
// 	categories := make([]Review, 0, len(dbCategories))
// 	for i, dbCat := range dbCategories {
// 		categories[i] = DatabaseCategoryToCategory(dbCat)
// 	}
// 	return categories
// }
