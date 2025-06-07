package handler

import (
	"backend/internal/common"
	"backend/internal/request"
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
	ctx := c.Request.Context()

	products, err := h.productService.GetAllProducts(ctx)
	if err != nil {
		fmt.Printf("Lỗi ở GetAllProductsService: %v\n", err)
		utils.JSON(c, http.StatusInternalServerError, "Không thể lấy danh sách sản phẩm", nil)
		return
	}

	utils.JSON(c, http.StatusOK, "lấy danh sách sản phẩm thành công", gin.H{
		"products": products,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	newProduct, err := h.productService.CreateProduct(ctx, req)
	if err != nil {
		fmt.Printf("Lỗi ở CreateProductService %v\n", err)
		utils.JSON(c, http.StatusInternalServerError, "Không thể tạo mới sản phẩm", nil)
	}

	utils.JSON(c, http.StatusOK, "Tạo mới sản phẩm thành công", gin.H{
		"product": newProduct,
	})
}
