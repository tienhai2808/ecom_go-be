package implement

import (
	customErr "backend/internal/errors"
	"backend/internal/model"
	"backend/internal/repository"
	"context"

	"gorm.io/gorm"
)

type addressRepositoryImpl struct {
	db *gorm.DB
}

func NewAddressRepository(db *gorm.DB) repository.AddressRepository {
	return &addressRepositoryImpl{
		db: db,
	}
}

func (r *addressRepositoryImpl) CreateAddress(ctx context.Context, address *model.Address) error {
	if err := r.db.WithContext(ctx).Create(address).Error; err != nil {
		return err
	}

	return nil
}

func (r *addressRepositoryImpl) GetAddressesByUserID(ctx context.Context, userID string) ([]*model.Address, error) {
	var addresses []*model.Address

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC").Find(&addresses).Error; err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *addressRepositoryImpl) CheckDefaultAddressExistsByUserID(ctx context.Context, userID string) (bool, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.Address{}).Where("user_id = ? AND is_default = true", userID).Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *addressRepositoryImpl) UpdateAddressIsDefaultByUserID(ctx context.Context, userID string, isDefault bool) error {
	result := r.db.WithContext(ctx).Model(&model.Address{}).Where("user_id = ? AND is_default = ?", userID, isDefault).Update("is_default", false)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrUserAddressNotFound
	}

	return nil
}
