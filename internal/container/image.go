package container

import (
	"backend/internal/handler"
	repoImpl "backend/internal/repository/implement"
	serviceImpl "backend/internal/service/implement"

	"gorm.io/gorm"
)

type ImageModule struct {
	ImageHandler handler.ImageHandler
}

func NewImageContainer(db *gorm.DB) *ImageModule {
	imageRepo := repoImpl.NewImageRepository(db)
	productRepo := repoImpl.NewProductRepository(db)
	imageService := serviceImpl.NewImageService(imageRepo, productRepo)
	imageHandler := handler.NewImageHandler(imageService)

	return &ImageModule{
		ImageHandler: *imageHandler,
	}
}