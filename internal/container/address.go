package container

import (
	"github.com/tienhai2808/ecom_go/internal/handler"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	svcImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"gorm.io/gorm"
)

type AddressModule struct {
	AddressHdl *handler.AddressHandler
}

func NewAddressContainer(db *gorm.DB, sfg snowflake.SnowflakeGenerator) *AddressModule {
	addressRepo := repoImpl.NewAddressRepository(db)
	addressSvc := svcImpl.NewAddressService(db, addressRepo, sfg)
	addressHdl := handler.NewAddressHandler(addressSvc)

	return &AddressModule{addressHdl}
}
