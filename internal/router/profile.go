package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"
)

func NewProfileRouter(rg *gin.RouterGroup, cfg *config.Config, userRepo repository.UserRepository, profileHdl *handler.ProfileHandler) {
	accessName := cfg.App.AccessName
	secretKey := cfg.App.JWTSecret

	profile := rg.Group("profiles")
	{
		profile.PATCH("/:id", security.RequireAuth(accessName, secretKey, userRepo), profileHdl.UpdateProfile)
	}
}
