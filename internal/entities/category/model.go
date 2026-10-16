package category

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	ThreadCount int
}

type CategoriesPage struct {
	Path       map[string]string
	Categories []Category
	Page       int
	Username   string
}
