package repository

import (
	"backend/internal/model"
	"context"
)

type AddressRepository interface {
	CreateAddress(ctx context.Context, address *model.Address) error

	GetAddressesByUserID(ctx context.Context, userID string) ([]*model.Address, error)

	CheckDefaultAddressExistsByUserID(ctx context.Context, userID string) (bool, error)

	UpdateAddressIsDefaultByUserID(ctx context.Context, userID string, isDefault bool) error 
}