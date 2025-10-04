package implement

import (
	"context"
	"errors"

	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"

	"gorm.io/gorm"
)

type profileRepositoryImpl struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) repository.ProfileRepository {
	return &profileRepositoryImpl{db}
}

func (r *profileRepositoryImpl) FindByID(ctx context.Context, id int64) (*model.Profile, error) {
	var profile model.Profile
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&profile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

func (r *profileRepositoryImpl) Update(ctx context.Context, id int64, updateData map[string]any) error {
	result := r.db.WithContext(ctx).Model(&model.Profile{}).Where("id = ?", id).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrProfileNotFound
	}

	return nil
}
