package router

import (
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(rg *gin.RouterGroup, config *config.Config, userRepository repository.UserRepository, userHandler *handler.UserHandler) {
	accessName := config.App.AccessName
	secretKey := config.App.JWTSecret

	user := rg.Group("/users", security.RequireAuth(accessName, secretKey, userRepository), security.RequireAdmin())
	{
		user.GET("", userHandler.GetAllUsers)

		user.GET("/:id", userHandler.GetUserByID)

		user.POST("", userHandler.CreateUser)

		user.PATCH("/:id", userHandler.UpdateUser)

		user.DELETE("/:id", userHandler.DeleteUser)

		user.DELETE("", userHandler.DeleteManyUsers)
	}
}
