package router

import (
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/common"
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
		address.GET("/my", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), addressHandler.GetMyAddresses)

		address.GET("/:id", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), addressHandler.GetAddressDetails)

		address.POST("", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), addressHandler.CreateAddress)

		address.PATCH("/:id", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), addressHandler.UpdateAddress)

		address.DELETE("/:id", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), addressHandler.DeleteAddress)
	}
}
