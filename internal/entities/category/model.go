package category

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type CategoriesPage struct {
	Path       string
	Categories []Category
	Page       int
	User       string
}
