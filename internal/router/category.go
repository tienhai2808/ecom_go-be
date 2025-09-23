package router

import (
	"backend/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func NewCategoryRouter(rg *gin.RouterGroup, config *config.Config, userRepository repository.UserRepository, categoryHandler *handler.CategoryHandler) {
	category := rg.Group("/categories")
	{
		category.POST("", middleware.RequireAuth(config, userRepository), middleware.RequireMultiRoles([]string{"admin"}), categoryHandler.CreateCategory)
	}
}