package implement

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/types"
	"gorm.io/gorm"
)

type cartRepositoryImpl struct {
	db  *gorm.DB
	rdb *redis.Client
	cfg *config.Config
}

func NewCartRepository(db *gorm.DB, rdb *redis.Client, cfg *config.Config) repository.CartRepository {
	return &cartRepositoryImpl{
		db,
		rdb,
		cfg,
	}
}

func (r *cartRepositoryImpl) GetGuestCartData(ctx context.Context, token string) (*types.CartData, error) {
	redisKey := fmt.Sprintf("%s:guest-cart:%s", r.cfg.App.Name, token)

	cartDataJSON, err := r.rdb.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("lấy dữ liệu từ redis thất bại: %w", err)
	}

	var cartData types.CartData
	if err = json.Unmarshal([]byte(cartDataJSON), &cartData); err != nil {
		return nil, fmt.Errorf("giải mã dữ liệu giỏ hàng thất bại: %w", err)
	}

	return &cartData, nil
}

func (r *cartRepositoryImpl) AddCartData(ctx context.Context, token string, data types.CartData, ttl time.Duration) error {
	cartDataJSON, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("mã hóa dữ liệu giỏ hàng thất bại: %w", err)
	}

	redisKey := fmt.Sprintf("%s:guest-cart:%s", r.cfg.App.Name, token)

	if err = r.rdb.Set(ctx, redisKey, cartDataJSON, ttl).Err(); err != nil {
		return err
	}

	return nil
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

func (r *cartRepositoryImpl) DeleteCartItemTx(ctx context.Context, tx *gorm.DB, cartItemID int64) error {
	return tx.WithContext(ctx).Where("id = ?", cartItemID).Delete(&model.CartItem{}).Error
}
