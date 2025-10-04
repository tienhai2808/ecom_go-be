package container

import (
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	serviceImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"gorm.io/gorm"
)

type ImageModule struct {
	ImageHandler *handler.ImageHandler
}

func NewImageContainer(db *gorm.DB, config *config.Config) *ImageModule {
	imageRepo := repoImpl.NewImageRepository(db)
	productRepo := repoImpl.NewProductRepository(db)
	imageService := serviceImpl.NewImageService(imageRepo, productRepo)
	imageHandler := handler.NewImageHandler(imageService)

	return &ImageModule{
		ImageHandler: imageHandler,
	}
}
