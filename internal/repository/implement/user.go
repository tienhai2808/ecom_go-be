package implement

import (
	customErr "backend/internal/errors"
	"backend/internal/model"
	"backend/internal/repository"
	"context"
	"errors"

	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) CheckUserExistsByEmail(ctx context.Context, email string) (bool, error) {
	var existingUser model.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First((&existingUser)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *userRepositoryImpl) CheckUserExistsByUsername(ctx context.Context, username string) (bool, error) {
	var existingUser model.User

	if err := r.db.WithContext(ctx).Where("username = ?", username).First((&existingUser)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *userRepositoryImpl) CreateUser(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepositoryImpl) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User

	if err := r.db.Preload("Profile").WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	if err := r.db.Preload("Profile").WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	if err := r.db.Preload("Profile").WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) UpdateUserPasswordByID(ctx context.Context, id, newPassword string) error {
	result := r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("password", newPassword)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return customErr.ErrUserNotFound
	}

	return nil
}
