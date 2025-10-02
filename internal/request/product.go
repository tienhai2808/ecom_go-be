package request

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,min=2"`
	CategoryID  int64   `json:"category_id" binding:"required,gt=0"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Quantity    uint    `json:"quantity" binding:"required,min=0"`
	Description string  `json:"description" binding:"omitempty"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=2"`
	CategoryID  *int64   `json:"category_id" binding:"omitempty,gt=0"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	Quantity    *uint    `json:"quantity" binding:"omitempty,min=0"`
	Description *string  `json:"description" binding:"omitempty"`
}
