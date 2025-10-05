package mapper

import (
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/response"
)

func ToCategoryResponse(category *model.Category) *response.CategoryResponse {
	return &response.CategoryResponse{
		ID: category.ID,
		Name: category.Name,
		Slug: category.Slug,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}