package implement

import (
	"context"
	"errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"

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

func (r *addressRepositoryImpl) FindByID(ctx context.Context, id string) (*model.Address, error) {
	return r.FindByIDTx(ctx, r.db, id)
}

func (r *addressRepositoryImpl) FindByIDTx(ctx context.Context, tx *gorm.DB, id string) (*model.Address, error) {
	var address model.Address
	if err := tx.WithContext(ctx).Where("id = ?", id).First(&address).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &address, nil
}

func (r *addressRepositoryImpl) Create(ctx context.Context, address *model.Address) error {
	return r.CreateTx(ctx, r.db, address)
}

func (r *addressRepositoryImpl) CreateTx(ctx context.Context, tx *gorm.DB, address *model.Address) error {
	return tx.WithContext(ctx).Create(address).Error
}

func (r *addressRepositoryImpl) FindByUserID(ctx context.Context, userID string) ([]*model.Address, error) {
	var addresses []*model.Address

	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("is_default DESC").Find(&addresses).Error; err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *addressRepositoryImpl) CountByUserIDTx(ctx context.Context, tx *gorm.DB, userID string) (int64, error) {
	var count int64
	if err := tx.WithContext(ctx).Model(&model.Address{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *addressRepositoryImpl) FindLatestByUserIDExcludeIDTx(ctx context.Context, tx *gorm.DB, userID, id string) (*model.Address, error) {
	var latestAddress model.Address
	if err := tx.WithContext(ctx).Where("user_id = ? AND id != ?", userID, id).Order("created_at DESC").First(&latestAddress).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &latestAddress, nil
}

func (r *addressRepositoryImpl) UpdateTx(ctx context.Context, tx *gorm.DB, id string, updateData map[string]any) error {
	return tx.WithContext(ctx).Model(&model.Address{}).Where("id = ?", id).Updates(updateData).Error
}

func (r *addressRepositoryImpl) DeleteTx(ctx context.Context, tx *gorm.DB, id string) error {
	return tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Address{}).Error
}
