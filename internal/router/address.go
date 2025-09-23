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
		address.GET("/my/all", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), addressHandler.GetUserAddresses)

		address.GET("/my/:address_id", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), addressHandler.GetUserAddressDetail)

		address.POST("/my", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), addressHandler.AddUserAddress)

		address.PATCH("/my/:address_id", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), addressHandler.UpdateUserAddress)

		address.DELETE("/my/:address_id", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), addressHandler.DeleteUserAddress)
	}
}
