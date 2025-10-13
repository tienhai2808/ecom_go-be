package repository

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"gorm.io/gorm"
)

type CartRepository interface {
	FindCartByUserID(ctx context.Context, userID int64) (*model.Cart, error)

	FindCartItemByCartIDAndProductIDTx(ctx context.Context, tx *gorm.DB, cartID, productID int64) (*model.CartItem, error)

	CreateCartItemTx(ctx context.Context, tx *gorm.DB, cartItem *model.CartItem) error

	CreateCart(ctx context.Context, cart *model.Cart) error

	UpdateCartTx(ctx context.Context, tx *gorm.DB, cartID int64, updateData map[string]any) error

	UpdateCartItemTx(ctx context.Context, tx *gorm.DB, cartItemID int64, updateData map[string]any) error

	FindCartByIDWithDetails(ctx context.Context, cartID int64) (*model.Cart, error)
}