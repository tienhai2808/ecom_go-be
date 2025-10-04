package handler

import (
	"github.com/tienhai2808/ecom_go/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categorySvc service.CategoryService
}

func NewCategoryHandler(categorySvc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categorySvc}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	c.JSON(200, "Hello")
}
