package service

import (
	"backend/internal/model"
	"backend/internal/request"
	"context"
)

type CategoryService interface {
	CreateCategory(ctx context.Context, req request.CreateCategoryRequest) (*model.Category, error)
}