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
	guestName := cfg.App.GuestName
	secretKey := cfg.App.JWTSecret

	cart := rg.Group("/carts", security.RequireAuth(accessName, secretKey, userRepo))
	{
		cart.POST("/items", cartHdl.AddCartItem)

		cart.PUT("/items/:id", cartHdl.UpdateCartItem)

		cart.GET("", cartHdl.GetMyCart)

		cart.DELETE("/items/:id", cartHdl.DeleteCartItem)
	}

	guest := rg.Group("/guests/carts", security.RequireGuestToken(guestName, secretKey))
	{
		guest.POST("/items", cartHdl.GuestAddCartItem)

		guest.GET("", cartHdl.GetGuestCart)
	}
}
