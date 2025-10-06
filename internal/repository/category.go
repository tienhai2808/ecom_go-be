package repository

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error

	FindAll(ctx context.Context) ([]*model.Category, error)
	
	FindByID(ctx context.Context, id int64) (*model.Category, error)
}