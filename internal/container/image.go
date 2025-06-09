package container

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/imagekit"
	repoImpl "backend/internal/repository/implement"
	serviceImpl "backend/internal/service/implement"

	"gorm.io/gorm"
)

type ImageModule struct {
	ImageHandler handler.ImageHandler
}

func NewImageContainer(db *gorm.DB, config *config.AppConfig) *ImageModule {
	imageRepo := repoImpl.NewImageRepository(db)
	productRepo := repoImpl.NewProductRepository(db)
	imageService := serviceImpl.NewImageService(imageRepo, productRepo)
	imageKitService := imagekit.NewImageKitService(config)
	imageHandler := handler.NewImageHandler(imageService, imageKitService)

	return &ImageModule{
		ImageHandler: *imageHandler,
	}
}