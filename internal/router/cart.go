package router

import (
	"github.com/gin-gonic/gin"
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"
)

func NewCartRouter(rg *gin.RouterGroup, cfg *config.Config, userRepo repository.UserRepository, cartHdl *handler.CartHandler) {
	accessName := cfg.App.AccessName
	secretKey := cfg.App.JWTSecret

	cart := rg.Group("/carts", security.RequireAuth(accessName, secretKey, userRepo))
	{
		cart.POST("/items", cartHdl.AddCartItem)
	}
}
