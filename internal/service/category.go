package service

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req request.CreateCategoryRequest) (*model.Category, error)

	GetAllCategories(ctx context.Context) ([]*model.Category, error)

	UpdateCategory(ctx context.Context, id int64, req request.UpdateCategoryRequest) (*model.Category, error)

	DeleteCategory(ctx context.Context, id int64) error

	DeleteCategories(ctx context.Context, req request.DeleteManyRequest) (int64, error)
}
