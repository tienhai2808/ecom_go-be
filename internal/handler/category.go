package handler

import (
	"github.com/tienhai2808/ecom_go/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	c.JSON(200, "Hello")
}
