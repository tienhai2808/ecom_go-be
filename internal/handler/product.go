package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/utils"

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
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	products, err := h.productService.GetAllProducts(ctx)
	if err != nil {
		utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "lấy danh sách sản phẩm thành công", gin.H{
		"products": products,
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	productID := c.Param("product_id")

	product, err := h.productService.GetProductByID(ctx, productID)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			utils.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Lấy sản phẩm thành công", gin.H{
		"product": product,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
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
		utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Tạo mới sản phẩm thành công", gin.H{
		"product": newProduct,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
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
			utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Cập nhật sản phẩm thành công", gin.H{
		"product": updatedProduct,
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	productID := c.Param("product_id")

	if err := h.productService.DeleteProduct(ctx, productID); err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			utils.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Xóa sản phẩm thành công", nil)
}

func (h *ProductHandler) DeleteManyProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.DeleteManyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	rowsAccepted, err := h.productService.DeleteManyProducts(ctx, req)
	if err != nil {
		utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	message := fmt.Sprintf("Xóa thành công %d người dùng", rowsAccepted)
	utils.JSON(c, http.StatusOK, message, nil)
}
