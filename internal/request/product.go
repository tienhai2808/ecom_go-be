package request

type UpdateProductForm struct {
	Name           *string                  `json:"name" validate:"omitempty,min=2"`
	CategoryID     *int64                   `json:"category_id" validate:"omitempty,gt=0"`
	Price          *float64                 `json:"price" validate:"omitempty,gt=0"`
	Quantity       *uint                    `json:"quantity" validate:"omitempty,min=0"`
	Description    *string                  `json:"description" validate:"omitempty,min=2"`
	IsActive       *bool                    `json:"is_active" validate:"omitempty"`
	NewImages      []CreateProductImageForm `json:"new_images" validate:"omitempty,dive"`
	UpdateImages   []UpdateProductImageForm `json:"update_images" validate:"omitempty,dive"`
	DeleteImageIDs []int64                  `json:"delete_image_ids" validate:"omitempty,dive"`
}

type UpdateProductImageForm struct {
	ID          int64 `json:"id" validate:"required"`
	IsThumbnail *bool `form:"is_thumbnail" validate:"omitempty"`
	SortOrder   *int  `form:"sort_order" validate:"omitempty,gt=0"`
}

type CreateProductForm struct {
	Name        string                   `json:"name" validate:"required,min=2"`
	CategoryID  int64                    `json:"category_id" validate:"required,gt=0"`
	Price       float64                  `json:"price" validate:"required,gt=0"`
	Quantity    uint                     `json:"quantity" validate:"required,min=1"`
	Description string                   `json:"description" validate:"required,min=2"`
	IsActive    *bool                    `json:"is_active" validate:"required"`
	Images      []CreateProductImageForm `json:"images" validate:"required,dive"`
}

type CreateProductImageForm struct {
	IsThumbnail *bool  `form:"is_thumbnail" validate:"required"`
	SortOrder   int    `form:"sort_order" validate:"required,gt=0"`
	FileData    []byte `form:"file_data" validate:"required"`
}

type ProductPaginationQuery struct {
	Page       uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit      uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort       string `form:"sort" json:"sort"`
	Order      string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	IsActive   *bool  `form:"is_active" json:"is_active"`
	Search     string `form:"search" json:"search"`
	CategoryID int64 `form:"category_id" json:"category_id" binding:"omitempty,gt=0"`
}
