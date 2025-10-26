package container

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	svcImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"gorm.io/gorm"
)

type CartModule struct {
	CartHdl *handler.CartHandler
}

func NewCartModule(db *gorm.DB, sfg snowflake.SnowflakeGenerator, es *elasticsearch.TypedClient, cfg *config.Config, rdb *redis.Client) *CartModule {
	cartRepo := repoImpl.NewCartRepository(db, rdb, cfg)
	productRepo := repoImpl.NewProductRepository(db, es)
	cartSvc := svcImpl.NewCartService(cartRepo, productRepo, db, sfg)
	cartHdl := handler.NewCartHandler(cartSvc)

	return &CartModule{cartHdl}
}
