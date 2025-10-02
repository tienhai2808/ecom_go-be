package container

import (
	"github.com/tienhai2808/ecom_go/internal/handler"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	serviceImpl "github.com/tienhai2808/ecom_go/internal/service/implement"

	"gorm.io/gorm"
)

type ProductModule struct {
	ProductHandler *handler.ProductHandler
}

func NewProductContainer(db *gorm.DB) *ProductModule {
	productRepo := repoImpl.NewProductRepository(db)
	categoryRepo := repoImpl.NewCategoryRepository(db)
	productService := serviceImpl.NewProductService(productRepo, categoryRepo)
	productHandler := handler.NewProductHandler(productService)

	return &ProductModule{
		ProductHandler: productHandler,
	}
}
