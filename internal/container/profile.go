package container

import (
	"github.com/tienhai2808/ecom_go/internal/handler"
	repoImpl "github.com/tienhai2808/ecom_go/internal/repository/implement"
	serviceImpl "github.com/tienhai2808/ecom_go/internal/service/implement"
	"gorm.io/gorm"
)

type ProfileModule struct {
	ProfileHdl *handler.ProfileHandler
}

func NewProfileContainer(db *gorm.DB) *ProfileModule {
	profileRepo := repoImpl.NewProfileRepository(db)
	userRepo := repoImpl.NewUserRepository(db)
	profileSvc := serviceImpl.NewProfileService(profileRepo, userRepo)
	profileHdl := handler.NewProfileHandler(profileSvc)

	return &ProfileModule{profileHdl}
}
