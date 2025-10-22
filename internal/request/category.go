package request

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=1"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=1"`
}
