package imagekit

import (
	"backend/internal/config"
	"mime/multipart"

	"github.com/imagekit-developer/imagekit-go"
)

type imageKitServicedImpl struct {
	client *imagekit.ImageKit
}

func NewImageKitService(config config.AppConfig) ImageKitService {
	params := imagekit.NewParams{
		PublicKey: config.ImageKit.PublicKey,
		PrivateKey: config.ImageKit.PrivateKey,
		UrlEndpoint: config.ImageKit.UrlEndpoint,
	}
	ik := imagekit.NewFromParams(params)

	return &imageKitServicedImpl{
		client: ik,
	}
}

func (s *imageKitServicedImpl) UploadImage(fileHeader *multipart.FileHeader) (string, string, error) {
	return "", "", nil
}
