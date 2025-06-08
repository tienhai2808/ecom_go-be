package handler

import (
	"backend/internal/common"
	customErr "backend/internal/errors"
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

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	ctx := c.Request.Context()
	productID := c.Param("product_id")

	product, err := h.productService.GetProductByID(ctx, productID)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			utils.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở GetProductByIDService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể lấy sản phẩm", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Lấy sản phẩm thành công", gin.H{
		"product": product,
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

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.UpdateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	productID := c.Param("product_id")

	updatedProduct, err := h.productService.UpdateProduct(ctx, productID, &req)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			utils.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở UpdateProductService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể cập nhât sản phẩm", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Cập nhật sản phẩm thành công", gin.H{
		"product": updatedProduct,
	})
}
