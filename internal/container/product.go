package container

import (
	"github.com/tienhai2808/ecom_go/internal/handler"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	serviceImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"github.com/tienhai2808/ecom_go/internal/snowflake"

	"gorm.io/gorm"
)

type ProductModule struct {
	ProductHandler *handler.ProductHandler
}

func NewProductContainer(db *gorm.DB, sfg snowflake.SnowflakeGenerator) *ProductModule {
	productRepo := repoImpl.NewProductRepository(db)
	categoryRepo := repoImpl.NewCategoryRepository(db)
	productSvc := serviceImpl.NewProductService(productRepo, categoryRepo, sfg)
	productHdl := handler.NewProductHandler(productSvc)

	return &ProductModule{productHdl}
}
