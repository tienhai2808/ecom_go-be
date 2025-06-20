package implement

import (
	customErr "backend/internal/errors"
	"backend/internal/model"
	"backend/internal/repository"
	"context"
	"errors"

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

func (r *addressRepositoryImpl) GetAddressByID(ctx context.Context, id string) (*model.Address, error) {
	var address model.Address

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &address, nil
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

func (r *addressRepositoryImpl) UpdateAddressNonDefaultByUserID(ctx context.Context, userID string) error {
	result := r.db.WithContext(ctx).Model(&model.Address{}).Where("user_id = ? AND is_default = true", userID).Update("is_default", false)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrUserAddressNotFound
	}

	return nil
}

func (r *addressRepositoryImpl) CountAddressByUserID(ctx context.Context, userID string) (int64, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&model.Address{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *addressRepositoryImpl) GetLatestAddressByUserIDExcludeID(ctx context.Context, userID, id string) (*model.Address, error) {
	var latestAddress model.Address

	if err := r.db.WithContext(ctx).Where("user_id = ? AND id != ?", userID, id).Order("created_at DESC").First(&latestAddress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &latestAddress, nil
}

func (r *addressRepositoryImpl) UpdateAddressByID(ctx context.Context, id string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Address{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrAddressNotFound
	}

	return nil
}

func (r *addressRepositoryImpl) UpdateAddressDefaultByID(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Model(&model.Address{}).Where("id = ?", id).Update("is_default", true)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrUserAddressNotFound
	}

	return nil
}

func (r *addressRepositoryImpl) DeleteAddressByID(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.Address{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrAddressNotFound
	}

	return nil
}
