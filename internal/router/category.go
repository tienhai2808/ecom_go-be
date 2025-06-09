package router

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func NewCategoryRouter(rg *gin.RouterGroup, config *config.AppConfig, userRepository repository.UserRepository, categoryHandler handler.CategoryHandler) {
	category := rg.Group("/categories")
	{
		category.POST("", middleware.RequireAuth(config, userRepository), middleware.RequireMultiRoles([]string{"admin"}), categoryHandler.CreateCategory)
	}
}