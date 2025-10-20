package implement

import (
	"context"
	"errors"

	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"gorm.io/gorm"
)

type cartRepositoryImpl struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) repository.CartRepository {
	return &cartRepositoryImpl{db}
}

func (r *cartRepositoryImpl) FindCartByUserID(ctx context.Context, userID int64) (*model.Cart, error) {
	var cart model.Cart
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepositoryImpl) FindCartByUserIDWithDetails(ctx context.Context, userID int64) (*model.Cart, error) {
	var cart model.Cart
	if err := r.db.WithContext(ctx).
		Preload("CartItems").
		Preload("CartItems.Product").
		Preload("CartItems.Product.Category").
		Preload("CartItems.Product.Images", "is_thumbnail = true").
		Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepositoryImpl) FindCartByUserIDTx(ctx context.Context, tx *gorm.DB, userID int64) (*model.Cart, error) {
	var cart model.Cart
	if err := tx.WithContext(ctx).Where("user_id = ?", userID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepositoryImpl) FindCartByIDWithDetails(ctx context.Context, cartID int64) (*model.Cart, error) {
	var cart model.Cart
	if err := r.db.WithContext(ctx).
		Preload("CartItems").
		Preload("CartItems.Product").
		Preload("CartItems.Product.Category").
		Preload("CartItems.Product.Images", "is_thumbnail = true").
		Where("id = ?", cartID).First(&cart).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &cart, nil
}

func (r *cartRepositoryImpl) FindCartItemByCartIDAndProductIDTx(ctx context.Context, tx *gorm.DB, cartID, productID int64) (*model.CartItem, error) {
	var cartItem model.CartItem
	if err := tx.WithContext(ctx).Where("cart_id = ? AND product_id = ?", cartID, productID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &cartItem, nil
}

func (r *cartRepositoryImpl) CreateCartItemTx(ctx context.Context, tx *gorm.DB, cartItem *model.CartItem) error {
	return tx.WithContext(ctx).Create(cartItem).Error
}

func (r *cartRepositoryImpl) CreateCart(ctx context.Context, cart *model.Cart) error {
	return r.db.WithContext(ctx).Create(cart).Error
}

func (r *cartRepositoryImpl) UpdateCartTx(ctx context.Context, tx *gorm.DB, cartID int64, updateData map[string]any) error {
	return tx.WithContext(ctx).Model(&model.Cart{}).Where("id = ?", cartID).Updates(updateData).Error
}

func (r *cartRepositoryImpl) UpdateCartItemTx(ctx context.Context, tx *gorm.DB, cartItemID int64, updateData map[string]any) error {
	return tx.WithContext(ctx).Model(&model.CartItem{}).Where("id = ?", cartItemID).Updates(updateData).Error
}

func (r *cartRepositoryImpl) FindCartItemByIDTx(ctx context.Context, tx *gorm.DB, cartItemID int64) (*model.CartItem, error) {
	var cartItem model.CartItem
	if err := r.db.WithContext(ctx).Where("id = ?", cartItemID).First(&cartItem).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &cartItem, nil
}
