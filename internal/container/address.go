package container

import (
	"github.com/tienhai2808/ecom_go/internal/handler"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	serviceImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"github.com/tienhai2808/ecom_go/internal/snowflake"
	"gorm.io/gorm"
)

type AddressModule struct {
	AddressHandler *handler.AddressHandler
}

func NewAddressContainer(db *gorm.DB, sfg snowflake.SnowflakeGenerator) *AddressModule {
	addressRepo := repoImpl.NewAddressRepository(db)
	addressService := serviceImpl.NewAddressService(db, addressRepo, sfg)
	addressHandler := handler.NewAddressHandler(addressService)

	return &AddressModule{
		AddressHandler: addressHandler,
	}
}
