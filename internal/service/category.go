package service

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req request.CreateCategoryRequest) (*model.Category, error)
}
