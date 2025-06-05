package repository

import (
	"backend/internal/model"
	"context"
)

type AddressRepository interface {
	GetAddressByID(ctx context.Context, id string) (*model.Address, error)

	CreateAddress(ctx context.Context, address *model.Address) error

	GetAddressesByUserID(ctx context.Context, userID string) ([]*model.Address, error)

	CheckDefaultAddressExistsByUserID(ctx context.Context, userID string) (bool, error)

	UpdateAddressNonDefaultByUserID(ctx context.Context, userID string) error 

	CountAddressByUserID(ctx context.Context, userID string) (int64, error)

	GetLatestAddressByUserIDExcludeID(ctx context.Context, userID, id string) (*model.Address, error)

	UpdateAddressByID(ctx context.Context, id string, updateData map[string]interface{}) error 

	UpdateAddressDefaultByID(ctx context.Context, id string) error

	DeleteAddressByID(ctx context.Context, id string) error
}