package response

type CartResponse struct {
	ID            int64               `json:"id"`
	TotalQuantity uint                `json:"total_quantity"`
	TotalPrice    float64             `json:"total_price"`
	CartItems     []*CartItemResponse `json:"cart_items"`
}

type CartItemResponse struct {
	ID         int64                  `json:"id"`
	UnitPrice  float64                `json:"unit_price"`
	Quantity   uint                   `json:"quantity"`
	TotalPrice float64                `json:"total_price"`
	Product    *SimpleProductResponse `json:"product"`
}
