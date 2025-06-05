package handler

import (
	"backend/internal/common"
	customErr "backend/internal/errors"
	"backend/internal/model"
	"backend/internal/request"
	"backend/internal/service"
	"backend/internal/utils"
	"fmt"
	"net/http"

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

func (h *AddressHandler) GetUserAddresses(c *gin.Context) {
	ctx := c.Request.Context()

	userAny, exists := c.Get("user")
	if !exists {
		utils.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		utils.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	addresses, err := h.addressService.GetUserAddresses(ctx, user.ID)
	if err != nil {
		fmt.Printf("Lỗi ở GetUserAddressService: %v\n", err)
		utils.JSON(c, http.StatusInternalServerError, "Không thể lấy địa chỉ của người dùng", nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Lấy địa chỉ người dùng thành công", gin.H{
		"addresses": addresses,
	})
}

func (h *AddressHandler) AddUserAddress(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.AddAddressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		utils.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		utils.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	newAddress, err := h.addressService.AddUserAddress(ctx, user.ID, req)
	if err != nil {
		switch err {
		case customErr.ErrUserAddressNotFound, customErr.ErrExceedsQuantity:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở AddUserAddressService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể thêm mới địa chỉ", nil)
		}
		return
	}

	utils.JSON(c, http.StatusCreated, "Thêm mới địa chỉ thành công", gin.H{
		"address": newAddress,
	})
}

func (h *AddressHandler) UpdateUserAddress(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.UpdateAddressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		utils.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		utils.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	addressID := c.Param("address_id")

	updatedAddress, err := h.addressService.UpdateUserAddress(ctx, user.ID, addressID, &req)
	if err != nil {
		switch err {
		case customErr.ErrAddressNotFound, customErr.ErrExceedsQuantity, customErr.ErrUserAddressNotFound, customErr.ErrUnauthorized:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở UpdateUserAddressService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể cập nhật địa chỉ người dùng", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Cập nhật địa chỉ người dùng thành công", gin.H{
		"address": updatedAddress,
	})
}

func (h *AddressHandler) DeleteUserAddress(c *gin.Context) {
	ctx := c.Request.Context()

	userAny, exists := c.Get("user")
	if !exists {
		utils.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		utils.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	addressID := c.Param("address_id")
	
	if err := h.addressService.DeleteUserAddress(ctx, user.ID, addressID); err != nil {
		switch err {
		case customErr.ErrAddressNotFound, customErr.ErrUnauthorized:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
	  default:
			fmt.Printf("Lỗi ở DeleteUserAddressService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể xóa địa chỉ", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Xóa địa chỉ thành công", nil)
}
