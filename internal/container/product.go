package container

import (
	"backend/internal/handler"
	repoImpl "backend/internal/repository/implement"
	serviceImpl "backend/internal/service/implement"

	"gorm.io/gorm"
)

type ProductModule struct {
	ProductHandler handler.ProductHandler
}

func NewProductContainer(db *gorm.DB) *ProductModule {
	productRepo := repoImpl.NewProductRepository(db)
	productService := serviceImpl.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	return &ProductModule{
		ProductHandler: *productHandler,
	}
}