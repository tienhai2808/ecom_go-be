package repository

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
)

type ImageRepository interface {
	GetImageByID(ctx context.Context, id string) (*model.Image, error)

	CreateImage(ctx context.Context, image *model.Image) error

	UpdateDownLoadUrlByID(ctx context.Context, id string, downloadUrl string) error
}
