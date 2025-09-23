package router

import (
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/common"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(rg *gin.RouterGroup, config *config.Config, userRepository repository.UserRepository, userHandler *handler.UserHandler) {
	accessName := config.App.AccessName
	secretKey := config.App.JWTSecret

	user := rg.Group("/users", security.RequireAuthAndRole(accessName, secretKey, common.RoleAdmin, userRepository))
	{
		user.GET("/all", userHandler.GetAllUsers)

		user.GET("/:user_id", userHandler.GetUserByID)

		user.POST("", userHandler.CreateUser)

		user.PATCH("/:user_id", userHandler.UpdateUser)

		user.DELETE("/:user_id", userHandler.DeleteUser)

		user.DELETE("/many", userHandler.DeleteManyUsers)
	}
}
