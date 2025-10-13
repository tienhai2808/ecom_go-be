package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/mapper"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/types"
)

type CartHandler struct {
	cartSvc service.CartService
}

func NewCartHandler(cartSvc service.CartService) *CartHandler {
	return &CartHandler{cartSvc}
}

func (h *CartHandler) AddCartItem(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		common.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	var req request.AddCartItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, translated, nil)
		return
	}

	cart, err := h.cartSvc.AddCartItem(ctx, user.ID, req)
	if err != nil {
		switch err {
		case customErr.ErrCartNotFound, customErr.ErrProductNotFound:
			common.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	common.JSON(c, http.StatusOK, "Thêm sản phẩm vào giỏ hàng thành công", gin.H{
		"cart": mapper.ToCartResponse(cart),
	})
}
