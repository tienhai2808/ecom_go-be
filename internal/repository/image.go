package repository

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
)

type ImageRepository interface {
	Create(ctx context.Context, image *model.Image) error

	Update(ctx context.Context, id int64, updateData map[string]any) error
}
