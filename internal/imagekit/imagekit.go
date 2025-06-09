package imagekit

import (
	"context"
	"mime/multipart"
)

type ImageKitService interface {
	UploadImage(ctx context.Context, fileName string, fileHeader *multipart.FileHeader) (string, string, error)
}