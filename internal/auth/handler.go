package auth

// import (
// 	"backend/internal/common"
// 	"fmt"
// 	"net/http"
// 	"time"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt/v5"
// )

// type Handler struct {
// 	service Service
// 	ctx     *common.AppContext
// }

// func NewHandler(service Service, ctx *common.AppContext) *Handler {
// 	return &Handler{
// 		service: service,
// 		ctx:     ctx,
// 	}
// }

// func (h *Handler) AddAddress(c *gin.Context) {
// 	var req AddAddressRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		translated := common.HandleValidationError(err)
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"errors":     translated,
// 			"statusCode": http.StatusBadRequest,
// 		})
// 		return
// 	}

// 	userIDVal, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"statusCode": http.StatusUnauthorized,
// 			"error":      "không có quyền truy cập",
// 		})
// 		return
// 	}

// 	userID, _ := userIDVal.(string)
// 	address, err := h.service.AddAddress(userID, req)
// 	if err != nil {
// 		switch err {
// 		case ErrUserNotFound:
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"statusCode": http.StatusBadRequest,
// 				"error":      err.Error(),
// 			})
// 		default:
// 			fmt.Printf("Lỗi ở AddAddressService: %v\n", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"statusCode": http.StatusInternalServerError,
// 				"error":      "Không thể thêm địa chỉ",
// 			})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusCreated, gin.H{
// 		"statusCode": http.StatusCreated,
// 		"address":    address,
// 	})
// }

// func (h *Handler) UpdateAddress(c *gin.Context) {
// 	var req UpdateAddressRequest

// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		translated := common.HandleValidationError(err)
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"errors":     translated,
// 			"statusCode": http.StatusBadRequest,
// 		})
// 		return
// 	}

// 	addressID := c.Param("address_id")

// 	userIDVal, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"statusCode": http.StatusUnauthorized,
// 			"error":      "không có quyền truy cập",
// 		})
// 		return
// 	}

// 	userID, _ := userIDVal.(string)

// 	updatedAddress, err := h.service.UpdateAddress(userID, addressID, req)
// 	if err != nil {
// 		switch err {
// 		case ErrUnAuth, ErrAddressNotFound, ErrUpdateFailed:
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"statusCode": http.StatusBadRequest,
// 				"error":      err.Error(),
// 			})
// 		default:
// 			fmt.Printf("Lỗi ở UpdateAddressService: %v\n", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"statusCode": http.StatusInternalServerError,
// 				"error":      "Không thể cập nhật địa chỉ",
// 			})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"statusCode": http.StatusOK,
// 		"address":    updatedAddress,
// 	})
// }

// func (h *Handler) GetAddresses(c *gin.Context) {
// 	userIDVal, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"statusCode": http.StatusUnauthorized,
// 			"error":      "không có quyền truy cập",
// 		})
// 		return
// 	}

// 	userID, _ := userIDVal.(string)

// 	addresses, err := h.service.GetAddresses(userID)
// 	if err != nil {
// 		fmt.Printf("Lỗi ở GetAddressesService: %v", err)
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"statusCode": http.StatusInternalServerError,
// 			"error":      "lỗi lấy dữ liệu địa chỉ người dùng",
// 		})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"statusCode": http.StatusOK,
// 		"addresses":  addresses,
// 	})
// }

// func (h *Handler) DeleteAddress(c *gin.Context) {
// 	userIDVal, exists := c.Get("user_id")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{
// 			"statusCode": http.StatusUnauthorized,
// 			"error":      "không có quyền truy cập",
// 		})
// 		return
// 	}

// 	userID, _ := userIDVal.(string)
// 	addressID := c.Param("address_id")
// 	if err := h.service.DeleteAddress(addressID, userID); err != nil {
// 		switch err {
// 		case ErrUnAuth, ErrAddressNotFound:
// 			c.JSON(http.StatusBadRequest, gin.H{
// 				"statusCode": http.StatusBadRequest,
// 				"error":      err.Error(),
// 			})
// 		default:
// 			fmt.Printf("Lỗi ở DeleteAddressService: %v\n", err)
// 			c.JSON(http.StatusInternalServerError, gin.H{
// 				"statusCode": http.StatusInternalServerError,
// 				"error":      "Không thể xóa địa chỉ",
// 			})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"statusCode": http.StatusOK,
// 		"message":    "Xóa địa chỉ thành công",
// 	})
// }
