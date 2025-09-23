package router

import (
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/common"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"

	"github.com/gin-gonic/gin"
)

func NewImageRouter(rg *gin.RouterGroup, config *config.Config, userRepository repository.UserRepository, imageHandler *handler.ImageHandler) {
	accessName := config.App.AccessName
	secretKey := config.App.JWTSecret

	image := rg.Group("/images")
	{
		image.POST("/upload", security.RequireAuthAndRole(accessName, secretKey, common.RoleAdmin, userRepository), imageHandler.UploadImages)
	}
}
