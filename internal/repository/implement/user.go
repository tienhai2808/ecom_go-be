package implement

import (
	customErr "backend/internal/errors"
	"backend/internal/model"
	"backend/internal/repository"

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

func (r *userRepositoryImpl) CheckUserExistsByEmail(email string) (bool, error) {
	var existingUser model.User

	err := r.db.Where("email = ?", email).First((&existingUser)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *userRepositoryImpl) CheckUserExistsByUsername(username string) (bool, error) {
	var existingUser model.User

	err := r.db.Where("username = ?", username).First((&existingUser)).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *userRepositoryImpl) CreateUser(user *model.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepositoryImpl) GetUserByUsername(username string) (*model.User, error) {
	var user model.User

	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, customErr.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}
