package service

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
	"mime/multipart"
)

type ImageService interface {
	UploadImages(ctx context.Context, files []*multipart.FileHeader, productID int64) ([]*model.Image, error)
}
