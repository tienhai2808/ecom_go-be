package container

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rabbitmq/amqp091-go"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	svcImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"gorm.io/gorm"
)

type ProductModule struct {
	ProductHdl *handler.ProductHandler
	ImageRepo  repository.ImageRepository
}

func NewProductContainer(db *gorm.DB, rabbitChan *amqp091.Channel, sfg snowflake.SnowflakeGenerator, es *elasticsearch.TypedClient) *ProductModule {
	productRepo := repoImpl.NewProductRepository(db, es)
	categoryRepo := repoImpl.NewCategoryRepository(db)
	inventoryRepo := repoImpl.NewInventoryRepository(db)
	imageRepo := repoImpl.NewImageRepository(db)
	productSvc := svcImpl.NewProductService(productRepo, categoryRepo, inventoryRepo, imageRepo, db, rabbitChan, sfg)
	productHdl := handler.NewProductHandler(productSvc)

	return &ProductModule{
		productHdl,
		imageRepo,
	}
}
