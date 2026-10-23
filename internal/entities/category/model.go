package category

import (
	"forum/internal/models"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ThreadCount int
}

type CategoriesPage struct {
	Breadcrumbs []models.Breadcrumb
	Categories  []Category
	Page        int
	Username    string
}
