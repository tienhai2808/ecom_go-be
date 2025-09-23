package router

import (
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/common"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"

	"github.com/gin-gonic/gin"
)

func NewCategoryRouter(rg *gin.RouterGroup, config *config.Config, userRepository repository.UserRepository, categoryHandler *handler.CategoryHandler) {
	accessName := config.App.AccessName
	secretKey := config.App.JWTSecret

	category := rg.Group("/categories")
	{
		category.POST("", security.RequireAuthAndRole(accessName, secretKey, common.RoleAdmin, userRepository), categoryHandler.CreateCategory)
	}
}
