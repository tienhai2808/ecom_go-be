package service

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
)

type CartService interface {
	AddCartItem(ctx context.Context, userID int64, req request.AddCartItemRequest) (*model.Cart, error)
}