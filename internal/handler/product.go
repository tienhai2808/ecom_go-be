package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
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
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "lấy danh sách sản phẩm thành công", gin.H{
		"products": products,
	})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		common.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	product, err := h.productSvc.GetProductByID(ctx, productID)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			common.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	common.JSON(c, http.StatusOK, "Lấy sản phẩm thành công", gin.H{
		"product": product,
	})
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		common.JSON(c, http.StatusBadRequest, customErr.ErrInvalidRequest.Error(), nil)
		return
	}

	var req request.CreateProductForm

	req.Name = strings.TrimSpace(c.PostForm("name"))

	if categoryIDStr := strings.TrimSpace(c.PostForm("category_id")); categoryIDStr != "" {
		if categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64); err == nil {
			req.CategoryID = categoryID
		}
	}

	if priceStr := strings.TrimSpace(c.PostForm("price")); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 32); err == nil {
			req.Price = price
		}
	}

	if quantityStr := strings.TrimSpace(c.PostForm("quantity")); quantityStr != "" {
		if quantity, err := strconv.Atoi(quantityStr); err == nil {
			req.Quantity = uint(quantity)
		}
	}

	req.Description = strings.TrimSpace(c.PostForm("description"))

	if isActiveStr := strings.TrimSpace(c.PostForm("is_active")); isActiveStr != "" {
		if isActive, err := strconv.ParseBool(isActiveStr); err != nil {
			req.IsActive = &isActive
		}
	}

	req.Images = []request.CreateProductImageForm{}
	i := 0
	for {
		isThumbnailKey := fmt.Sprintf("images[%d][is_thumbnail]", i)
		sortOrderKey := fmt.Sprintf("images[%d][sort_order]", i)
		fileKey := fmt.Sprintf("images[%d][file]", i)

		isThumbnailStr := strings.TrimSpace(c.PostForm(isThumbnailKey))
		if isThumbnailStr == "" {
			break
		}

		isThumbnail := false
		if isThumbnailStr != "" {
			isThumbnail, _ = strconv.ParseBool(isThumbnailStr)
		}

		sortOrder := 0
		sortOrderStr := strings.TrimSpace(c.PostForm(sortOrderKey))
		if sortOrderStr != "" {
			sortOrder, _ = strconv.Atoi(sortOrderStr)
		}

		file, err := c.FormFile(fileKey)
		if err != nil {
			common.JSON(c, http.StatusBadRequest, fmt.Sprintf("Không tìm thấy file cho ảnh %d", i), nil)
			return
		}

		openedFile, err := file.Open()
		if err != nil {
			common.JSON(c, http.StatusBadRequest, fmt.Sprintf("Không mở được file cho ảnh %d", i), nil)
			return
		}
		defer openedFile.Close()

		fileBytes, err := io.ReadAll(openedFile)
		if err != nil {
			common.JSON(c, http.StatusBadRequest, fmt.Sprintf("Đọc file cho ảnh %d thất bại", i), nil)
			return
		}

		image := request.CreateProductImageForm{
			IsThumbnail: &isThumbnail,
			SortOrder:   sortOrder,
			File:        fileBytes,
		}

		req.Images = append(req.Images, image)
		i++
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		translated := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, translated, nil)
		return
	}

	newProduct, err := h.productSvc.CreateProduct(ctx, &req)
	if err != nil {
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Tạo mới sản phẩm thành công", gin.H{
		"product": newProduct,
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, customErr.ErrInvalidRequest.Error(), gin.H{
			"errors": translated,
		})
		return
	}

	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		common.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	updatedProduct, err := h.productSvc.UpdateProduct(ctx, productID, &req)
	if err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			common.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	common.JSON(c, http.StatusOK, "Cập nhật sản phẩm thành công", gin.H{
		"product": updatedProduct,
	})
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	productIDStr := c.Param("id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		common.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	if err := h.productSvc.DeleteProduct(ctx, productID); err != nil {
		switch err {
		case customErr.ErrProductNotFound:
			common.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	common.JSON(c, http.StatusOK, "Xóa sản phẩm thành công", nil)
}

func (h *ProductHandler) DeleteManyProducts(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.DeleteManyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, customErr.ErrInvalidRequest.Error(), gin.H{
			"errors": translated,
		})
		return
	}

	rowsAccepted, err := h.productSvc.DeleteManyProducts(ctx, req)
	if err != nil {
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	message := fmt.Sprintf("Xóa thành công %d người dùng", rowsAccepted)
	common.JSON(c, http.StatusOK, message, nil)
}
