package container

import (
	"github.com/tienhai2808/ecom_go/internal/handler"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	serviceImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"gorm.io/gorm"
)

type CategoryModule struct {
	CategoryHdl *handler.CategoryHandler
}

func NewCategoryContainer(db *gorm.DB, sfg snowflake.SnowflakeGenerator) *CategoryModule {
	categoryRepo := repoImpl.NewCategoryRepository(db)
	categorySvc := serviceImpl.NewCategoryService(categoryRepo, sfg)
	categoryHdl := handler.NewCategoryHandler(categorySvc)

	return &CategoryModule{categoryHdl}
}
