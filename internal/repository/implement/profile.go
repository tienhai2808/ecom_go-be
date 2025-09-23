package implement

import (
	"context"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/repository"

	"gorm.io/gorm"
)

type profileRepositoryImpl struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) repository.ProfileRepository {
	return &profileRepositoryImpl{
		db: db,
	}
}

func (r *profileRepositoryImpl) UpdateProfileByUserID(ctx context.Context, userID string, updateData map[string]interface{}) error {
	result := r.db.WithContext(ctx).Model(&model.Profile{}).Where("user_id = ?", userID).Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrUserProfileNotFound
	}

	return nil
}
