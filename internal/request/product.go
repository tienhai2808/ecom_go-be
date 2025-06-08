package request

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=2"`
	Brand       string  `json:"brand" binding:"required,min=2,max=50"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Inventory   uint    `json:"inventory" binding:"required,min=0"`
	Description string  `json:"description" binding:"omitempty"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=2"`
	Brand       *string  `json:"brand" binding:"omitempty,min=2,max=50"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	Inventory   *uint    `json:"inventory" binding:"omitempty,min=0"`
	Description *string  `json:"description" binding:"omitempty"`
}
