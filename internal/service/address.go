package service

import (
	"context"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
)

type AddressService interface {
	GetMyAddresses(ctx context.Context, userID string) ([]*model.Address, error)

	GetAddressDetail(ctx context.Context, userID string, id string) (*model.Address, error)

	CreateAddress(ctx context.Context, userID string, req request.AddAddressRequest) (*model.Address, error)

	UpdateAddress(ctx context.Context, userID, id string, req *request.UpdateAddressRequest) (*model.Address, error)

	DeleteAddress(ctx context.Context, userID, id string) error
}
