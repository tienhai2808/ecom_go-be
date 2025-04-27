package admin

import (
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
			"error": "lỗi lấy dữ liệu người dùng",
		})
		return 
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"users": users,
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
		fmt.Printf("Lỗi ở CreateUsersService: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError, 
			"error": "lỗi tạo mới người dùng",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": http.StatusCreated,
		"message": "Tạo mới người dùng thành công",
	})
}