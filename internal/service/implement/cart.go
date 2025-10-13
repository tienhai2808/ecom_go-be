package implement

import (
	"context"
	"fmt"

	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"gorm.io/gorm"
)

type cartServiceImpl struct {
	cartRepo    repository.CartRepository
	productRepo repository.ProductRepository
	db          *gorm.DB
	sfg         snowflake.SnowflakeGenerator
}

func NewCartService(cartRepo repository.CartRepository, productRepo repository.ProductRepository, db *gorm.DB, sfg snowflake.SnowflakeGenerator) service.CartService {
	return &cartServiceImpl{
		cartRepo,
		productRepo,
		db,
		sfg,
	}
}

func (s *cartServiceImpl) AddCartItem(ctx context.Context, userID int64, req request.AddCartItemRequest) (*model.Cart, error) {
	cart, err := s.cartRepo.FindCartByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin giỏ hàng thất bại: %w", err)
	}
	if cart == nil {
		cartID, err := s.sfg.NextID()
		if err != nil {
			return nil, err
		}

		cart = &model.Cart{
			ID:            cartID,
			TotalPrice:    0,
			TotalQuantity: 0,
			UserID:        userID,
		}
		if err = s.cartRepo.CreateCart(ctx, cart); err != nil {
			return nil, fmt.Errorf("tạo giỏ hàng thất bại: %w", err)
		}
	}

	product, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
	}
	if product == nil {
		return nil, customErr.ErrProductNotFound
	}

	if err = s.db.Transaction(func(tx *gorm.DB) error {
		existingItem, err := s.cartRepo.FindCartItemByCartIDAndProductIDTx(ctx, tx, cart.ID, product.ID)
		if err != nil {
			return fmt.Errorf("kiểm tra sản phẩm trong giỏ hàng thất bại: %w", err)
		}

		if existingItem != nil {
			newQuantity := existingItem.Quantity + req.Quantity
			newTotalPrice := float64(newQuantity) * existingItem.UnitPrice

			updateData := map[string]any{
				"total_price": newTotalPrice,
				"quantity":    newQuantity,
			}
			if err = s.cartRepo.UpdateCartItemTx(ctx, tx, existingItem.ID, updateData); err != nil {
				return fmt.Errorf("cập nhật sản phẩm trong giỏ hàng thất bại: %w", err)
			}

			totalPriceCart := cart.TotalPrice + product.Price*float64(req.Quantity)
			totalQuantityCart := cart.TotalQuantity + req.Quantity

			updateData = map[string]any{
				"total_price":    totalPriceCart,
				"total_quantity": totalQuantityCart,
			}
			if err = s.cartRepo.UpdateCartTx(ctx, tx, cart.ID, updateData); err != nil {
				return fmt.Errorf("cập nhật giỏ hàng thất bại: %w", err)
			}
		} else {
			cartItemID, err := s.sfg.NextID()
			if err != nil {
				return err
			}

			cartItem := &model.CartItem{
				ID:        cartItemID,
				UnitPrice: product.Price,
				Quantity:  req.Quantity,
				CartID:    cart.ID,
				ProductID: product.ID,
			}
			cartItem.SetTotalPrice()

			if err = s.cartRepo.CreateCartItemTx(ctx, tx, cartItem); err != nil {
				return fmt.Errorf("thêm sản phẩm vào giỏ hàng thất bại: %w", err)
			}

			totalPriceCart := cart.TotalPrice + cartItem.TotalPrice
			totalQuantityCart := cart.TotalQuantity + cartItem.Quantity

			updateData := map[string]any{
				"total_price":    totalPriceCart,
				"total_quantity": totalQuantityCart,
			}
			if err = s.cartRepo.UpdateCartTx(ctx, tx, cart.ID, updateData); err != nil {
				return fmt.Errorf("cập nhật giỏ hàng thất bại: %w", err)
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	cart, err = s.cartRepo.FindCartByIDWithDetails(ctx, cart.ID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin giỏ hàng thất bại: %w", err)
	}
	if cart == nil {
		return nil, customErr.ErrCartNotFound
	}

	return cart, nil
}
