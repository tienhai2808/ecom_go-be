package cloudinary

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService interface {
	UploadBinaryFile(ctx context.Context, fileData []byte, fileName string) (*uploader.UploadResult, error)

	DeleteFile(ctx context.Context, fileName, fileType string) error
}
