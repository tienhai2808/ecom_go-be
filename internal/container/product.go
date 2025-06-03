package container

import (
	"backend/internal/handler"
	repository "backend/internal/repository/implement"
	service "backend/internal/service/implement"

	"gorm.io/gorm"
)

type ProductModule struct {
	ProductHandler handler.ProductHandler
}

func NewProductContainer(db *gorm.DB) *ProductModule {
	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	return &ProductModule{
		ProductHandler: *productHandler,
	}
}