package handler

import (
	"backend/internal/service"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	products, err := h.productService.GetAllProducts()
	if err != nil {
		fmt.Println("Lỗi ở lấy tất cả sản phẩm: ", err)
		utils.JSON(c, http.StatusInternalServerError, "lỗi lấy dữ liệu sản phẩm", nil)
		return 
	}

	utils.JSON(c, http.StatusOK, "lấy tất cả sản phẩm thành công", products)
}