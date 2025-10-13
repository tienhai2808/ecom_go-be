package mapper

import (
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/response"
)

func ToCartResponse(cart *model.Cart) *response.CartResponse {
	return &response.CartResponse{
		ID: cart.ID,
		TotalQuantity: cart.TotalQuantity,
		TotalPrice: cart.TotalPrice,
		CartItems: ToCartItemsResponse(cart.CartItems),
	}
}

func ToCartItemResponse(cartItem *model.CartItem) *response.CartItemResponse {
	return &response.CartItemResponse{
		ID: cartItem.ID,
		UnitPrice: cartItem.UnitPrice,
		Quantity: cartItem.Quantity,
		TotalPrice: cartItem.TotalPrice,
		Product: ToSimpleProductResponse(cartItem.Product),
	}
}

func ToCartItemsResponse(cIts []*model.CartItem) []*response.CartItemResponse {
	if len(cIts) == 0 {
		return make([]*response.CartItemResponse, 0)
	}

	cItsResp := make([]*response.CartItemResponse, 0, len(cIts))
	for _, cIt := range cIts {
		cItsResp = append(cItsResp, ToCartItemResponse(cIt))
	}

	return cItsResp
}