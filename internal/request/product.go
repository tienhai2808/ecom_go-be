package request

type CreateProductForm struct {
	Name        string                   `json:"name" validate:"required,min=2"`
	CategoryID  int64                    `json:"category_id" validate:"required,gt=0"`
	Price       float64                  `json:"price" validate:"required,gt=0"`
	Quantity    uint                     `json:"quantity" validate:"required,min=0"`
	Description string                   `json:"description" validate:"required,min=2"`
	IsActive    *bool                    `json:"is_active" validate:"required"`
	Images      []CreateProductImageForm `json:"images" validate:"required,dive"`
}

type CreateProductImageForm struct {
	IsThumbnail *bool  `form:"is_thumbnail" validate:"required"`
	SortOrder   int    `form:"sort_order" validate:"required,gt=0"`
	FileData    []byte `form:"file_data" validate:"required"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name" binding:"omitempty,min=2"`
	CategoryID  *int64   `json:"category_id" binding:"omitempty,gt=0"`
	Price       *float64 `json:"price" binding:"omitempty,gt=0"`
	Quantity    *uint    `json:"quantity" binding:"omitempty,min=0"`
	Description *string  `json:"description" binding:"omitempty"`
}
