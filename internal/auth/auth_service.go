package auth

import (
	"e-commerce/internal/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SignupService(db *gorm.DB, dto SignupSchema) (*user.User, error) {
	hashedPassword, err := hashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	newUser := &user.User{
		ID:       uuid.NewString(),
		Username: dto.Username,
		Email:    dto.Email,
		Password: hashedPassword,
		Profile: user.Profile{
			ID: uuid.NewString(),
		},
	}

	if err := db.Create(newUser).Error; err != nil {
		return nil, err
	}

	return newUser, nil
}
