package service

import (
	"backend/internal/model"
	"context"
	"mime/multipart"
)

type ImageService interface {
	UploadImages(ctx context.Context, files []*multipart.FileHeader, productID string) ([]*model.Image, error)
}