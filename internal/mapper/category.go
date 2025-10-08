package mapper

import (
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/response"
)

func ToCategoryResponse(category *model.Category) *response.CategoryResponse {
	return &response.CategoryResponse{
		ID:        category.ID,
		Name:      category.Name,
		Slug:      category.Slug,
		CreatedAt: category.CreatedAt,
		UpdatedAt: category.UpdatedAt,
	}
}

func ToBaseCategoryResponse(category *model.Category) *response.BaseCategoryResponse {
	return &response.BaseCategoryResponse{
		ID:   category.ID,
		Name: category.Name,
		Slug: category.Slug,
	}
}

func ToCategoriesResponse(ctgs []*model.Category) []*response.BaseCategoryResponse {
	if len(ctgs) == 0 {
		return make([]*response.BaseCategoryResponse, 0)
	}

	ctgsResp := make([]*response.BaseCategoryResponse, 0, len(ctgs))
	for _, ctg := range ctgs {
		ctgsResp = append(ctgsResp, ToBaseCategoryResponse(ctg))
	}

	return ctgsResp
}
