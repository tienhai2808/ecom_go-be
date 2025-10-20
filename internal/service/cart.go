package service

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
)

type CartService interface {
	GetMyCart(ctx context.Context, userID int64) (*model.Cart, error)

	AddCartItem(ctx context.Context, userID int64, req request.AddCartItemRequest) (*model.Cart, error)
  
	UpdateCartItem(ctx context.Context, userID, cartItemID int64, quantity uint) (*model.Cart, error)

	DeleteCartItem(ctx context.Context, userID, cartItemID int64) (*model.Cart, error)
}