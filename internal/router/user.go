package router

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(rg *gin.RouterGroup, config *config.AppConfig, userRepository repository.UserRepository, userHandler handler.UserHandler) {
	user := rg.Group("/users", middleware.RequireAuth(config, userRepository), middleware.RequireMultiRoles([]string{"admin"}))
	{
		user.GET("/all", userHandler.GetAllUsers)

		user.POST("", userHandler.CreateUser)

		user.PATCH("/:user_id", userHandler.UpdateUser)

		user.DELETE("/many", userHandler.DeleteManyUsers)

		// address.POST("/my", middleware.RequireAuth(config, userRepository), addressHandler.AddUserAddress)

		// address.PATCH("/my/:address_id", middleware.RequireAuth(config, userRepository), addressHandler.UpdateUserAddress)

		// address.DELETE("/my/:address_id", middleware.RequireAuth(config, userRepository), addressHandler.DeleteUserAddress)
	}
}