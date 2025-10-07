package repository

import (
	"context"

	"github.com/tienhai2808/ecom_go/internal/model"
	"gorm.io/gorm"
)

type ImageRepository interface {
	Create(ctx context.Context, image *model.Image) error

	CreateAllTx(ctx context.Context, tx *gorm.DB, images []*model.Image) error

	Update(ctx context.Context, id int64, updateData map[string]any) error

	UpdateTx(ctx context.Context, tx *gorm.DB, id int64, updateData map[string]any) error

	FindAllByIDTx(ctx context.Context, tx *gorm.DB, ids []int64) ([]*model.Image, error)

	DeleteAllByIDTx(ctx context.Context, tx *gorm.DB, ids []int64) error
}
