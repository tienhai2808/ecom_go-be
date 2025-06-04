package implement

import (
	customErr "backend/internal/errors"
	"backend/internal/model"
	"backend/internal/repository"
	"context"

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
