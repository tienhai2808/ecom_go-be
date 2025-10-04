package router

import (
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"

	"github.com/gin-gonic/gin"
)

func NewImageRouter(rg *gin.RouterGroup, cfg *config.Config, userRepo repository.UserRepository, imageHdl *handler.ImageHandler) {
	accessName := cfg.App.AccessName
	secretKey := cfg.App.JWTSecret

	image := rg.Group("/images")
	{
		image.POST("/upload", security.RequireAuth(accessName, secretKey, userRepo), security.RequireAdmin(), imageHdl.UploadImages)
	}
}
