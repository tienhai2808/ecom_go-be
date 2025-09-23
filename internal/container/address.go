package container

import (
	"backend/internal/handler"
	serviceImpl "backend/internal/service/implement"
	repoImpl "backend/internal/repository/implement"

	"gorm.io/gorm"
)

type AddressModule struct {
	AddressHandler *handler.AddressHandler
}

func NewAddressContainer(db *gorm.DB) *AddressModule {
	addressRepo := repoImpl.NewAddressRepository(db)
	addressService := serviceImpl.NewAddressService(addressRepo)
	addressHandler := handler.NewAddressHandler(addressService)

	return &AddressModule{
		AddressHandler: addressHandler,
	}
}