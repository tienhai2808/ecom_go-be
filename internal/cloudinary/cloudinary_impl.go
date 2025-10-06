package cloudinary

import (
	"bytes"
	"context"
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type cloudinaryServiceImpl struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService(cld *cloudinary.Cloudinary) CloudinaryService {
	return &cloudinaryServiceImpl{cld}
}

func (s *cloudinaryServiceImpl) UploadBinaryFile(ctx context.Context, fileData []byte, fileName string) (*uploader.UploadResult, error) {
	res, err := s.cld.Upload.Upload(ctx, bytes.NewReader(fileData), uploader.UploadParams{
		Folder:         "ecom_go/product",
		UniqueFilename: toBoolPnt(true),
		PublicID:       fileName,
		Overwrite:      toBoolPnt(false),
	})
	if err != nil {
		return nil, fmt.Errorf("đăng tải binary file thất bại: %w", err)
	}

	return res, nil
}

func (s *cloudinaryServiceImpl) DeleteFile(ctx context.Context, fileName, fileType string) error {
	if _, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     fileName,
		ResourceType: fileType,
	}); err != nil {
		return fmt.Errorf("xóa file thất bại: %w", err)
	}

	return nil
}

func toBoolPnt(b bool) *bool {
	return &b
}
