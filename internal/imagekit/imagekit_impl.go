package imagekit

import (
	"backend/config"
	"context"
	"fmt"
	"mime/multipart"

	"github.com/imagekit-developer/imagekit-go"
	"github.com/imagekit-developer/imagekit-go/api/uploader"
)

type imageKitServicedImpl struct {
	Client *imagekit.ImageKit
}

func NewImageKitService(config *config.Config) ImageKitService {
	params := imagekit.NewParams{
		UrlEndpoint: config.ImageKit.UrlEndpoint,
		PublicKey:   config.ImageKit.PublicKey,
		PrivateKey:  config.ImageKit.PrivateKey,
	}
	ik := imagekit.NewFromParams(params)
	return &imageKitServicedImpl{
		Client: ik,
	}
}

func (s *imageKitServicedImpl) UploadImage(ctx context.Context, fileName string, fileHeader *multipart.FileHeader) (string, string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", "", fmt.Errorf("mở file thất bại: %w", err)
	}
	defer file.Close()

	res, err := s.Client.Uploader.Upload(ctx, file, uploader.UploadParam{
		FileName: fileName,
		Folder:   "ecom-go",
	})

	if err != nil {
		return "", "", fmt.Errorf("upload ảnh thất bại: %w", err)
	}

	return res.Data.Url, res.Data.FileId, nil
}
