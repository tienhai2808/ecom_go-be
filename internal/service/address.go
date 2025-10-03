package service

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
)

type AddressService interface {
	GetMyAddresses(ctx context.Context, userID int64) ([]*model.Address, error)

	GetAddressDetail(ctx context.Context, userID int64, id int64) (*model.Address, error)

	CreateAddress(ctx context.Context, userID int64, req request.AddAddressRequest) (*model.Address, error)

	UpdateAddress(ctx context.Context, userID, id int64, req request.UpdateAddressRequest) (*model.Address, error)

	DeleteAddress(ctx context.Context, userID, id int64) error
}
