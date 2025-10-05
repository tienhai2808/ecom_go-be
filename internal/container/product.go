package container

import (
	"github.com/rabbitmq/amqp091-go"
	"github.com/tienhai2808/ecom_go/internal/handler"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	serviceImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"github.com/tienhai2808/ecom_go/internal/snowflake"

	"gorm.io/gorm"
)

type ProductModule struct {
	ProductHdl *handler.ProductHandler
}

func NewProductContainer(db *gorm.DB, rabbitChan *amqp091.Channel, sfg snowflake.SnowflakeGenerator) *ProductModule {
	productRepo := repoImpl.NewProductRepository(db)
	categoryRepo := repoImpl.NewCategoryRepository(db)
	productSvc := serviceImpl.NewProductService(productRepo, categoryRepo, rabbitChan, sfg)
	productHdl := handler.NewProductHandler(productSvc)

	return &ProductModule{productHdl}
}
