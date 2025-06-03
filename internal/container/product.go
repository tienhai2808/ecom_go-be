package container

import (
	"backend/internal/handler"
	"backend/internal/implement"

	"gorm.io/gorm"
)

type ProductModule struct {
	ProductHandler handler.ProductHandler
}

func NewProductContainer(db *gorm.DB) *ProductModule {
	repo := implement.NewProductRepository(db)
	service := implement.NewProductService(repo)
	handler := handler.NewProductHandler(service)

	return &ProductModule{
		ProductHandler: *handler,
	}
}