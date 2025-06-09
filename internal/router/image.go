package router

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func NewImageRouter(rg *gin.RouterGroup, config *config.AppConfig, userRepository repository.UserRepository, imageHandler handler.ImageHandler) {
	image := rg.Group("/images")
	{
		image.POST("/upload", middleware.RequireAuth(config, userRepository), middleware.RequireMultiRoles([]string{"admin"}), imageHandler.UploadImages)
	}
}