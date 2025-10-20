package request

type AddCartItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required,gt=0"`
	Quantity  uint  `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemRequest struct {
	Quantity uint `json:"quantity" binding:"required,min=1"`
}