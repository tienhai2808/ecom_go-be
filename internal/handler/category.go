package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/mapper"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
)

type CategoryHandler struct {
	categorySvc service.CategoryService
}

func NewCategoryHandler(categorySvc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categorySvc}
}

func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	var req request.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, translated, nil)
		return
	}

	category, err := h.categorySvc.CreateCategory(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrCategorySlugAlreadyExists:
			common.JSON(c, http.StatusConflict, err.Error(), nil)
		default:
			common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	common.JSON(c, http.StatusCreated, "Tạo danh mục sản phẩm thành công", gin.H{
		"category": mapper.ToCategoryResponse(category),
	})
}

func (h *CategoryHandler) GetAllCategories(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	categories, err := h.categorySvc.GetAllCategories(ctx)
	if err != nil {
		common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	common.JSON(c, http.StatusOK, "Lấy danh sách danh mục sản phẩm thành công", gin.H{
		"categories": mapper.ToCategoriesResponse(categories),
	})
}

func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	categoryIDStr := c.Param("id")
	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		common.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	var req request.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, translated, nil)
		return
	}

	category, err := h.categorySvc.UpdateCategory(ctx, categoryID, req)
	if err != nil {
		switch err {
		case customErr.ErrCategoryNotFound:
			common.JSON(c, http.StatusNotFound, err.Error(), nil)
		case customErr.ErrCategorySlugAlreadyExists:
			common.JSON(c, http.StatusConflict, err.Error(), nil)
		default:
			common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	common.JSON(c, http.StatusOK, "Cập nhật danh mục sản phẩm thành công", gin.H{
		"category": mapper.ToCategoryResponse(category),
	})
}
