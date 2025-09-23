package container

import (
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	serviceImpl "github.com/tienhai2808/ecom_go/internal/service/implement"

	"gorm.io/gorm"
)

type UserModule struct {
	UserRepository repository.UserRepository
	UserHandler    *handler.UserHandler
}

func NewUserContainer(db *gorm.DB) *UserModule {
	userRepo := repoImpl.NewUserRepository(db)
	profileRepo := repoImpl.NewProfileRepository(db)
	userService := serviceImpl.NewUserService(userRepo, profileRepo)
	userHandler := handler.NewUserHandler(userService)

	return &UserModule{
		UserRepository: userRepo,
		UserHandler:    userHandler,
	}
}
