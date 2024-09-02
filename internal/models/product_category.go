package models

import (
	"github.com/CP-Payne/ecomstore/internal/database"
	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
}

// Database Category to Category mappings
func DatabaseCategoryToCategory(category database.Category) Category {
	return Category{
		ID:          category.ID,
		Name:        category.Name,
		Description: NullStringToString(category.Description),
	}
}

func DatabaseCategoriesToCategories(dbCategories []database.Category) []Category {
	categories := make([]Category, len(dbCategories))
	for i, dbCat := range dbCategories {
		categories[i] = DatabaseCategoryToCategory(dbCat)
	}
	return categories
}
