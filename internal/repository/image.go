package repository

import (
	"backend/internal/model"
	"context"
)

type ImageRepository interface {
	GetImageByID(ctx context.Context, id string) (*model.Image, error)

	CreateImage(ctx context.Context, image *model.Image) error

	UpdateDownLoadUrlByID(ctx context.Context, id string, downloadUrl string) error
}