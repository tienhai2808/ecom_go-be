package router

import (
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"
	"github.com/gin-gonic/gin"
)

func NewAddressRouter(rg *gin.RouterGroup, cfg *config.Config, userRepo repository.UserRepository, addressHdl *handler.AddressHandler) {
	accessName := cfg.App.AccessName
	secretKey := cfg.App.JWTSecret

	address := rg.Group("/addresses")
	{
		address.GET("/my", security.RequireAuth(accessName, secretKey, userRepo), addressHdl.GetMyAddresses)

		address.GET("/:id", security.RequireAuth(accessName, secretKey, userRepo), addressHdl.GetAddressDetails)

		address.POST("", security.RequireAuth(accessName, secretKey, userRepo), addressHdl.CreateAddress)

		address.PATCH("/:id", security.RequireAuth(accessName, secretKey, userRepo), addressHdl.UpdateAddress)

		address.DELETE("/:id", security.RequireAuth(accessName, secretKey, userRepo), addressHdl.DeleteAddress)
	}
}
