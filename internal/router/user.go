package router

import (
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"

	"github.com/gin-gonic/gin"
)

func NewUserRouter(rg *gin.RouterGroup, cfg *config.Config, userRepo repository.UserRepository, userHdl *handler.UserHandler) {
	accessName := cfg.App.AccessName
	secretKey := cfg.App.JWTSecret

	user := rg.Group("/users", security.RequireAuth(accessName, secretKey, userRepo), security.RequireAdmin())
	{
		user.GET("", userHdl.GetAllUsers)

		user.GET("/:id", userHdl.GetUserByID)

		user.POST("", userHdl.CreateUser)

		user.PATCH("/:id", userHdl.UpdateUser)

		user.DELETE("/:id", userHdl.DeleteUser)

		user.DELETE("", userHdl.DeleteManyUsers)
	}
}
