package auth

import (
	"e-commerce/internal/user"
	"errors"

	"gorm.io/gorm"
)

func IsUserExists(db *gorm.DB, username, email string) error {
	var existingUser user.User

	if err := db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err 
	}

	if existingUser.Username == username {
		return errors.New("username already taken")
	}
	if existingUser.Email == email {
		return errors.New("email already registered")
	}

	return nil
}
