package admin

import (
	"backend/internal/auth"
	"backend/internal/common"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
	ctx     *common.AppContext
}

func NewHandler(service Service, ctx *common.AppContext) *Handler {
	return &Handler{
		service: service,
		ctx:     ctx,
	}
}

func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		fmt.Printf("Lỗi ở GetAllUsersService: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error":      "lỗi lấy dữ liệu người dùng",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"users":      users,
	})
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	if err := h.service.CreateUser(req); err != nil {
		switch err {
		case auth.ErrUsernameExists, auth.ErrEmailExists:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		default:
			fmt.Printf("Lỗi ở SignupService: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error":      "Không thể đăng ký tài khoản",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"message":    "Tạo mới người dùng thành công",
	})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	userID := c.Param("user_id")
	if err := h.service.UpdateUser(userID, &req); err != nil {
		switch err {
		case auth.ErrUsernameExists, auth.ErrEmailExists, auth.ErrUpdateFailed, auth.ErrUserNotFound:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		default:
			fmt.Printf("Lỗi ở SignupService: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error":      "Không thể cập nhật người dùng",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Cập nhật thông tin người dùng thành công",
	})
}
