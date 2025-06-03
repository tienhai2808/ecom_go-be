package repository

import "backend/internal/model"

type UserRepository interface {
	CheckUserExistsByEmail(email string) (bool, error)
	CheckUserExistsByUsername(username string) (bool, error)
	CreateUser(user *model.User) error
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(userID string) (*model.User, error)
}
