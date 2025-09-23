package router

import (
	"backend/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(rg *gin.RouterGroup, config *config.Config, userRepository repository.UserRepository, userHandler *handler.UserHandler) {
	user := rg.Group("/users", middleware.RequireAuth(config, userRepository), middleware.RequireMultiRoles([]string{"admin"}))
	{
		user.GET("/all", userHandler.GetAllUsers)

		user.GET("/:user_id", userHandler.GetUserByID)

		user.POST("", userHandler.CreateUser)

		user.PATCH("/:user_id", userHandler.UpdateUser)

		user.DELETE("/:user_id", userHandler.DeleteUser)

		user.DELETE("/many", userHandler.DeleteManyUsers)
	}
}