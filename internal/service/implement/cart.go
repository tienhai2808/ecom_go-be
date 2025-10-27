package implement

import (
	"context"
	"fmt"
	"time"

	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/mapper"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/response"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"github.com/tienhai2808/ecom_go/internal/types"
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

func (s *cartServiceImpl) GetMyCart(ctx context.Context, userID int64) (*model.Cart, error) {
	cart, err := s.cartRepo.FindCartByUserIDWithDetails(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin giỏ hàng thất bại: %w", err)
	}
	if cart == nil {
		return nil, customErr.ErrCartNotFound
	}

	return cart, nil
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

func (s *cartServiceImpl) UpdateCartItem(ctx context.Context, userID, cartItemID int64, quantity uint) (*model.Cart, error) {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		cart, err := s.cartRepo.FindCartByUserIDTx(ctx, tx, userID)
		if err != nil {
			return fmt.Errorf("lấy thông tin giỏ hàng thất bại: %w", err)
		}
		if cart == nil {
			return customErr.ErrCartNotFound
		}

		cartItem, err := s.cartRepo.FindCartItemByIDTx(ctx, tx, cartItemID)
		if err != nil {
			return fmt.Errorf("lấy thông tin mặt hàng trong giỏ hàng thất bại: %w", err)
		}
		if cartItem == nil {
			return customErr.ErrCartItemNotFound
		}

		if cartItem.CartID != cart.ID {
			return customErr.ErrCartItemNotFound
		}

		newTotalPrice := cartItem.UnitPrice * float64(quantity)
		updateData := map[string]any{
			"quantity":    quantity,
			"total_price": newTotalPrice,
		}
		if err := s.cartRepo.UpdateCartItemTx(ctx, tx, cartItemID, updateData); err != nil {
			return fmt.Errorf("cập nhật mặt hàng trong giỏ hàng thất bại: %w", err)
		}

		totalPriceCart := cart.TotalPrice - cartItem.TotalPrice + cartItem.UnitPrice*float64(quantity)
		totalQuantityCart := cart.TotalQuantity - cartItem.Quantity + quantity
		updateData = map[string]any{
			"total_price":    totalPriceCart,
			"total_quantity": totalQuantityCart,
		}
		if err = s.cartRepo.UpdateCartTx(ctx, tx, cart.ID, updateData); err != nil {
			return fmt.Errorf("cập nhật giỏ hàng thất bại: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	updatedCart, err := s.cartRepo.FindCartByUserIDWithDetails(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin giỏ hàng thất bại: %w", err)
	}
	if updatedCart == nil {
		return nil, customErr.ErrCartNotFound
	}

	return updatedCart, nil
}

func (s *cartServiceImpl) DeleteCartItem(ctx context.Context, userID, cartItemID int64) (*model.Cart, error) {
	if err := s.db.Transaction(func(tx *gorm.DB) error {
		cart, err := s.cartRepo.FindCartByUserIDTx(ctx, tx, userID)
		if err != nil {
			return fmt.Errorf("lấy thông tin giỏ hàng thất bại: %w", err)
		}
		if cart == nil {
			return customErr.ErrCartNotFound
		}

		cartItem, err := s.cartRepo.FindCartItemByIDTx(ctx, tx, cartItemID)
		if err != nil {
			return fmt.Errorf("lấy thông tin mặt hàng trong giỏ hàng thất bại: %w", err)
		}
		if cartItem == nil {
			return customErr.ErrCartItemNotFound
		}

		if cartItem.CartID != cart.ID {
			return customErr.ErrCartItemNotFound
		}

		if err = s.cartRepo.DeleteCartItemTx(ctx, tx, cartItemID); err != nil {
			return fmt.Errorf("xóa mặt hàng khỏi giỏ hàng thất bại: %w", err)
		}

		totalPriceCart := cart.TotalPrice - cartItem.TotalPrice
		totalQuantityCart := cart.TotalQuantity - cartItem.Quantity
		updateData := map[string]any{
			"total_price":    totalPriceCart,
			"total_quantity": totalQuantityCart,
		}
		if err = s.cartRepo.UpdateCartTx(ctx, tx, cart.ID, updateData); err != nil {
			return fmt.Errorf("cập nhật giỏ hàng thất bại: %w", err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	updatedCart, err := s.cartRepo.FindCartByUserIDWithDetails(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin giỏ hàng thất bại: %w", err)
	}
	if updatedCart == nil {
		return nil, customErr.ErrCartNotFound
	}

	return updatedCart, nil
}

func (s *cartServiceImpl) GuestAddCartItem(ctx context.Context, guestID string, req request.AddCartItemRequest) (*response.GuestCartResponse, error) {
	cart, err := s.cartRepo.GetGuestCartData(ctx, guestID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		cart = &types.CartData{
			TotalPrice:    0,
			TotalQuantity: 0,
			Items:         []types.CartItemData{},
		}
	}

	product, err := s.productRepo.FindByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("lấy thông tin sản phẩm thất bại: %w", err)
	}
	if product == nil {
		return nil, customErr.ErrProductNotFound
	}

	found := false
	for i := range cart.Items {
		if cart.Items[i].ProductID == product.ID {
			cart.Items[i].Quantity += req.Quantity
			cart.Items[i].TotalPrice = cart.Items[i].UnitPrice * float64(cart.Items[i].Quantity)

			cart.TotalQuantity += req.Quantity
			cart.TotalPrice += cart.Items[i].UnitPrice * float64(req.Quantity)

			found = true
			break
		}
	}

	if !found {
		newItem := types.CartItemData{
			ProductID:  product.ID,
			Quantity:   req.Quantity,
			UnitPrice:  product.Price,
			TotalPrice: product.Price * float64(req.Quantity),
		}
		cart.Items = append(cart.Items, newItem)

		cart.TotalQuantity += req.Quantity
		cart.TotalPrice += newItem.UnitPrice * float64(newItem.Quantity)
	}

	if err = s.cartRepo.AddCartData(ctx, guestID, *cart, 7*24*time.Hour); err != nil {
		return nil, err
	}

	productIDs := make([]int64, 0, len(cart.Items))
	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := s.productRepo.FindAllByIDWithCategoryAndThumbnail(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("lây thông tin sản phẩm trong giỏ hàng thất bại: %w", err)
	}
	if len(productIDs) != len(products) {
		return nil, customErr.ErrHasProductNotFound
	}

	productMap := make(map[int64]*model.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	return toGuestCartResponse(cart, productMap), nil
}

func (s *cartServiceImpl) GetGuestCart(ctx context.Context, guestID string) (*response.GuestCartResponse, error) {
	cart, err := s.cartRepo.GetGuestCartData(ctx, guestID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		cart = &types.CartData{
			TotalPrice:    0,
			TotalQuantity: 0,
			Items:         []types.CartItemData{},
		}

		if err = s.cartRepo.AddCartData(ctx, guestID, *cart, 7*24*time.Hour); err != nil {
			return nil, err
		}
	}

	if len(cart.Items) == 0 {
		return &response.GuestCartResponse{
			TotalQuantity: cart.TotalQuantity,
			TotalPrice:    cart.TotalPrice,
			CartItems:     []*response.GuestCartItemResponse{},
		}, nil
	}

	productIDs := make([]int64, 0, len(cart.Items))
	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := s.productRepo.FindAllByIDWithCategoryAndThumbnail(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("lây thông tin sản phẩm trong giỏ hàng thất bại: %w", err)
	}
	if len(productIDs) != len(products) {
		return nil, customErr.ErrHasProductNotFound
	}

	productMap := make(map[int64]*model.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	return toGuestCartResponse(cart, productMap), nil
}

func (s *cartServiceImpl) GuestUpdateCartItem(ctx context.Context, guestID string, productID int64, quantity uint) (*response.GuestCartResponse, error) {
	cart, err := s.cartRepo.GetGuestCartData(ctx, guestID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, customErr.ErrCartNotFound
	}

	found := false
	for i := range cart.Items {
		if cart.Items[i].ProductID == productID {
			oldQty := cart.Items[i].Quantity
			oldTotal := cart.Items[i].TotalPrice

			cart.Items[i].Quantity = quantity
			cart.Items[i].TotalPrice = cart.Items[i].UnitPrice * float64(quantity)

			diffQty := int(quantity) - int(oldQty)
			diffTotal := cart.Items[i].TotalPrice - oldTotal

			newTotalQty := int(cart.TotalQuantity) + diffQty
			cart.TotalQuantity = uint(newTotalQty)
			cart.TotalPrice += diffTotal

			found = true
			break
		}
	}

	if !found {
		return nil, customErr.ErrCartItemNotFound
	}

	if err = s.cartRepo.AddCartData(ctx, guestID, *cart, 7*24*time.Hour); err != nil {
		return nil, err
	}

	productIDs := make([]int64, 0, len(cart.Items))
	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := s.productRepo.FindAllByIDWithCategoryAndThumbnail(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("lây thông tin sản phẩm trong giỏ hàng thất bại: %w", err)
	}
	if len(productIDs) != len(products) {
		return nil, customErr.ErrHasProductNotFound
	}

	productMap := make(map[int64]*model.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	return toGuestCartResponse(cart, productMap), nil
}

func (s *cartServiceImpl) GuestDeleteCartItem(ctx context.Context, guestID string, productID int64) (*response.GuestCartResponse, error) {
	cart, err := s.cartRepo.GetGuestCartData(ctx, guestID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, customErr.ErrCartNotFound
	}

	found := false
	newItems := make([]types.CartItemData, 0, len(cart.Items))
	for _, item := range cart.Items {
		if item.ProductID == productID {
			cart.TotalQuantity -= item.Quantity
			cart.TotalPrice -= item.TotalPrice

			found = true
			continue
		}
		newItems = append(newItems, item)
	}

	if !found {
		return nil, customErr.ErrCartItemNotFound
	}

	cart.Items = newItems

	if err = s.cartRepo.AddCartData(ctx, guestID, *cart, 7*24*time.Hour); err != nil {
		return nil, err
	}

	productIDs := make([]int64, 0, len(cart.Items))
	for _, item := range cart.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	products, err := s.productRepo.FindAllByIDWithCategoryAndThumbnail(ctx, productIDs)
	if err != nil {
		return nil, fmt.Errorf("lây thông tin sản phẩm trong giỏ hàng thất bại: %w", err)
	}
	if len(productIDs) != len(products) {
		return nil, customErr.ErrHasProductNotFound
	}

	productMap := make(map[int64]*model.Product, len(products))
	for _, p := range products {
		productMap[p.ID] = p
	}

	return toGuestCartResponse(cart, productMap), nil
}

func toGuestCartResponse(cart *types.CartData, productMap map[int64]*model.Product) *response.GuestCartResponse {
	cartItemsResp := make([]*response.GuestCartItemResponse, 0, len(cart.Items))
	for _, item := range cart.Items {
		p := productMap[item.ProductID]

		var prodResp *response.SimpleProductResponse
		if p != nil {
			prodResp = mapper.ToSimpleProductResponse(p)
		}

		cartItemsResp = append(cartItemsResp, &response.GuestCartItemResponse{
			UnitPrice:  item.UnitPrice,
			Quantity:   item.Quantity,
			TotalPrice: item.TotalPrice,
			Product:    prodResp,
		})
	}

	return &response.GuestCartResponse{
		TotalQuantity: cart.TotalQuantity,
		TotalPrice:    cart.TotalPrice,
		CartItems:     cartItemsResp,
	}
}
