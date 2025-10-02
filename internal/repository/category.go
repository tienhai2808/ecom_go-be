package repository

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
)

type CategoryRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Category, error)
}