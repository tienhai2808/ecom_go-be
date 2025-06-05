package service

import (
	"backend/internal/model"
	"backend/internal/request"
	"context"
)

type AddressService interface {
	GetUserAddresses(ctx context.Context, userID string) ([]*model.Address, error)

	AddUserAddress(ctx context.Context, userID string, req request.AddAddressRequest) (*model.Address, error)

	UpdateUserAddress(ctx context.Context, userID, id string, req *request.UpdateAddressRequest) (*model.Address, error)

	DeleteUserAddress(ctx context.Context, userID, id string) error
}
