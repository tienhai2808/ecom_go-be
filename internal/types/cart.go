package types

type CartData struct {
	TotalPrice    float64        `json:"total_price"`
	TotalQuantity uint           `json:"total_quantity"`
	Items         []CartItemData `json:"items"`
}

type CartItemData struct {
	UnitPrice  float64 `json:"unit_price"`
	Quantity   uint    `json:"quantity"`
	TotalPrice float64 `json:"total_price"`
	ProductID  int64   `json:"product_id"`
}
