package repository

import (
	"context"
	"time"

	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/types"
	"gorm.io/gorm"
)

type CartRepository interface {
	FindCartByUserID(ctx context.Context, userID int64) (*model.Cart, error)

	FindCartByUserIDWithDetails(ctx context.Context, userID int64) (*model.Cart, error)

	FindCartByUserIDTx(ctx context.Context, tx *gorm.DB, userID int64) (*model.Cart, error)

	FindCartItemByCartIDAndProductIDTx(ctx context.Context, tx *gorm.DB, cartID, productID int64) (*model.CartItem, error)

	CreateCartItemTx(ctx context.Context, tx *gorm.DB, cartItem *model.CartItem) error

	CreateCart(ctx context.Context, cart *model.Cart) error

	UpdateCartTx(ctx context.Context, tx *gorm.DB, cartID int64, updateData map[string]any) error

	UpdateCartItemTx(ctx context.Context, tx *gorm.DB, cartItemID int64, updateData map[string]any) error

	FindCartByIDWithDetails(ctx context.Context, cartID int64) (*model.Cart, error)

	FindCartItemByIDTx(ctx context.Context, tx *gorm.DB, cartItemID int64) (*model.CartItem, error)

	DeleteCartItemTx(ctx context.Context, tx *gorm.DB, cartItemID int64) error

	GetGuestCartData(ctx context.Context, token string) (*types.CartData, error)

	AddCartData(ctx context.Context, token string, data types.CartData, ttl time.Duration) error
}