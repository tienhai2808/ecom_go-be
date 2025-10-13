package container

import (
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	svcImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"github.com/tienhai2808/ecom_go/internal/snowflake"

	"gorm.io/gorm"
)

type UserModule struct {
	UserRepo repository.UserRepository
	UserHdl  *handler.UserHandler
}

func NewUserContainer(db *gorm.DB, sfg snowflake.SnowflakeGenerator) *UserModule {
	userRepo := repoImpl.NewUserRepository(db)
	profileRepo := repoImpl.NewProfileRepository(db)
	userSvc := svcImpl.NewUserService(userRepo, profileRepo, sfg)
	userHdl := handler.NewUserHandler(userSvc)

	return &UserModule{
		userRepo,
		userHdl,
	}
}
