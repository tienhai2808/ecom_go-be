package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"
)

func NewProfileRouter(rg *gin.RouterGroup, config *config.Config, userRepository repository.UserRepository, handler *handler.ProfileHandler) {
	accessName := config.App.AccessName
	secretKey := config.App.JWTSecret

	profile := rg.Group("profiles")
	{
		profile.PATCH("/:id", security.RequireAuth(accessName, secretKey, userRepository), handler.UpdateProfile)
	}
}
