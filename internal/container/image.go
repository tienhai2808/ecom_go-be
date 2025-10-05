package container

import (
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	serviceImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"gorm.io/gorm"
)

type ImageModule struct {
	ImageHdl  *handler.ImageHandler
	ImageRepo repository.ImageRepository
}

func NewImageContainer(db *gorm.DB, config *config.Config) *ImageModule {
	imageRepo := repoImpl.NewImageRepository(db)
	productRepo := repoImpl.NewProductRepository(db)
	imageSvc := serviceImpl.NewImageService(imageRepo, productRepo)
	imageHdl := handler.NewImageHandler(imageSvc)

	return &ImageModule{
		imageHdl,
		imageRepo,
	}
}
