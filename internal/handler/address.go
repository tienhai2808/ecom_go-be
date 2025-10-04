package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/mapper"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/util"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	addressService service.AddressService
}

func NewAddressHandler(addressService service.AddressService) *AddressHandler {
	return &AddressHandler{
		addressService: addressService,
	}
}

func (h *AddressHandler) GetMyAddresses(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		util.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		util.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	addresses, err := h.addressService.GetMyAddresses(ctx, user.ID)
	if err != nil {
		util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	util.JSON(c, http.StatusOK, "Lấy địa chỉ người dùng thành công", gin.H{
		"addresses": mapper.ToAddressesResponse(addresses),
	})
}

func (h *AddressHandler) GetAddressDetails(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		util.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		util.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	addressIDStr := c.Param("id")
	addressID, err := strconv.ParseInt(addressIDStr, 10, 64)
	if err != nil {
		util.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	address, err := h.addressService.GetAddressDetail(ctx, user.ID, addressID)
	if err != nil {
		switch err {
		case customErr.ErrUnauthorized, customErr.ErrAddressNotFound:
			util.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	util.JSON(c, http.StatusOK, "Lấy địa chỉ thành công", gin.H{
		"address": mapper.ToAddressResponse(address),
	})
}

func (h *AddressHandler) CreateAddress(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.AddAddressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		util.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		util.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		util.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	newAddress, err := h.addressService.CreateAddress(ctx, user.ID, req)
	if err != nil {
		switch err {
		case customErr.ErrUserAddressNotFound, customErr.ErrExceedsQuantity:
			util.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	util.JSON(c, http.StatusCreated, "Thêm mới địa chỉ thành công", gin.H{
		"address": mapper.ToAddressResponse(newAddress),
	})
}

func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.UpdateAddressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		util.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		util.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		util.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	addressIDStr := c.Param("id")
	addressID, err := strconv.ParseInt(addressIDStr, 10, 64)
	if err != nil {
		util.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	updatedAddress, err := h.addressService.UpdateAddress(ctx, user.ID, addressID, req)
	if err != nil {
		switch err {
		case customErr.ErrAddressNotFound, customErr.ErrExceedsQuantity, customErr.ErrUserAddressNotFound, customErr.ErrUnauthorized:
			util.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	util.JSON(c, http.StatusOK, "Cập nhật địa chỉ người dùng thành công", gin.H{
		"address": mapper.ToAddressResponse(updatedAddress),
	})
}

func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userAny, exists := c.Get("user")
	if !exists {
		util.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		util.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	addressIDStr := c.Param("id")
	addressID, err := strconv.ParseInt(addressIDStr, 10, 64)
	if err != nil {
		util.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	if err := h.addressService.DeleteAddress(ctx, user.ID, addressID); err != nil {
		switch err {
		case customErr.ErrAddressNotFound, customErr.ErrUnauthorized:
			util.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	util.JSON(c, http.StatusOK, "Xóa địa chỉ thành công", nil)
}
