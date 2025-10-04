package router

import (
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"

	"github.com/gin-gonic/gin"
)

func NewAddressRouter(rg *gin.RouterGroup, config *config.Config, userRepository repository.UserRepository, addressHandler *handler.AddressHandler) {
	accessName := config.App.AccessName
	secretKey := config.App.JWTSecret

	address := rg.Group("/addresses")
	{
		address.GET("/my", security.RequireAuth(accessName, secretKey, userRepository), addressHandler.GetMyAddresses)

		address.GET("/:id", security.RequireAuth(accessName, secretKey, userRepository), addressHandler.GetAddressDetails)

		address.POST("", security.RequireAuth(accessName, secretKey, userRepository), addressHandler.CreateAddress)

		address.PATCH("/:id", security.RequireAuth(accessName, secretKey, userRepository), addressHandler.UpdateAddress)

		address.DELETE("/:id", security.RequireAuth(accessName, secretKey, userRepository), addressHandler.DeleteAddress)
	}
}
