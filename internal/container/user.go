package container

import (
	"backend/internal/handler"
	"backend/internal/repository"
	repoImpl "backend/internal/repository/implement"
	serviceImpl "backend/internal/service/implement"

	"gorm.io/gorm"
)

type UserModule struct {
	UserRepository repository.UserRepository
	UserHandler    handler.UserHandler
}

func NewUserContainer(db *gorm.DB) *UserModule {
	userRepo := repoImpl.NewUserRepository(db)
	profileRepo := repoImpl.NewProfileRepository(db)
	userService := serviceImpl.NewUserService(userRepo, profileRepo)
	userHandler := handler.NewUserHandler(userService)

	return &UserModule{
		UserRepository: userRepo,
		UserHandler: *userHandler,
	}
}