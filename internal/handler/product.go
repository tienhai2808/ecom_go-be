package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/util"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productSvc service.ProductService
}

func NewProductHandler(productSvc service.ProductService) *ProductHandler {
	return &ProductHandler{productSvc}
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	products, err := h.productSvc.GetAllProducts(ctx)
	if err != nil {
		util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	util.JSON(c, http.StatusOK, "lấy danh sách sản phẩm thành công", gin.H{
		"products": products,
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		util.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	product, err := h.productSvc.GetProductByID(ctx, productID)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			util.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	util.JSON(c, http.StatusOK, "Lấy sản phẩm thành công", gin.H{
		"product": product,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		util.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	newProduct, err := h.productSvc.CreateProduct(ctx, req)
	if err != nil {
		util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	util.JSON(c, http.StatusOK, "Tạo mới sản phẩm thành công", gin.H{
		"product": newProduct,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.UpdateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		util.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		util.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	updatedProduct, err := h.productSvc.UpdateProduct(ctx, productID, &req)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			util.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	util.JSON(c, http.StatusOK, "Cập nhật sản phẩm thành công", gin.H{
		"product": updatedProduct,
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		util.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	if err := h.productSvc.DeleteProduct(ctx, productID); err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			util.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	util.JSON(c, http.StatusOK, "Xóa sản phẩm thành công", nil)
}

func (h *ProductHandler) DeleteManyProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.DeleteManyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		util.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	rowsAccepted, err := h.productSvc.DeleteManyProducts(ctx, req)
	if err != nil {
		util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	message := fmt.Sprintf("Xóa thành công %d người dùng", rowsAccepted)
	util.JSON(c, http.StatusOK, message, nil)
}
