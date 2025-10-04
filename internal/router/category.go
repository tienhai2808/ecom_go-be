package router

import (
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"
	"github.com/gin-gonic/gin"
)

func NewCategoryRouter(rg *gin.RouterGroup, cfg *config.Config, userRepo repository.UserRepository, categoryHdl *handler.CategoryHandler) {
	accessName := cfg.App.AccessName
	secretKey := cfg.App.JWTSecret

	category := rg.Group("/categories")
	{
		category.POST("", security.RequireAuth(accessName, secretKey, userRepo), security.RequireAdmin(), categoryHdl.CreateCategory)
	}
}
