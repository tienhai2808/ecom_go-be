package router

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func NewAddressRouter(rg *gin.RouterGroup, config *config.AppConfig, userRepository repository.UserRepository, addressHandler handler.AddressHandler) {
	address := rg.Group("/addresses")
	{
		address.GET("/my", middleware.RequireAuth(config, userRepository), addressHandler.GetUserAddresses)

		address.POST("/my", middleware.RequireAuth(config, userRepository), addressHandler.AddUserAddress)

		address.PATCH("/my/:address_id", middleware.RequireAuth(config, userRepository), addressHandler.UpdateUserAddress)

		address.DELETE("/my/:address_id", middleware.RequireAuth(config, userRepository), addressHandler.DeleteUserAddress)
	}
}
