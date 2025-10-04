package initialization

import (
	"fmt"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/tienhai2808/ecom_go/config"
)

func InitCloudinary(cfg *config.Config) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(cfg.Cloudinary.CloudName, cfg.Cloudinary.ApiKey, cfg.Cloudinary.ApiSecret)
	if err != nil {
		return nil, fmt.Errorf("khởi tạo Cloudinary thất bại: %w", err)
	}
	cld.Config.URL.Secure = true

	return cld, nil
}
