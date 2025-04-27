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
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"users": users,
	})
}
