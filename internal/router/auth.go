package router

import (
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/common"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(rg *gin.RouterGroup, config *config.Config, userRepository repository.UserRepository, authHandler *handler.AuthHandler) {
	accessName := config.App.AccessName
	refreshName := config.App.RefreshName
	secretKey := config.App.JWTSecret

	auth := rg.Group("/auth")
	{
		auth.POST("/signup", authHandler.Signup)

		auth.POST("/signup/verify", authHandler.VerifySignup)

		auth.POST("/signin", authHandler.Signin)

		auth.POST("/signout", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), authHandler.Signout)

		auth.GET("/me", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), authHandler.GetMe)

		auth.GET("/refresh-token", security.RequireRefreshToken(refreshName, secretKey, userRepository), authHandler.RefreshToken)

		auth.POST("/forgot-password", authHandler.ForgotPassword)

		auth.POST("/forgot-password/verify", authHandler.VerifyForgotPassword)

		auth.POST("/reset-password", authHandler.ResetPassword)

		auth.PUT("/change-password/:user_id", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), authHandler.ChangePassword)

		auth.PATCH("/update-info/:user_id", security.RequireAuthAndRole(accessName, secretKey, common.RoleUser, userRepository), authHandler.UpdateUserProfile)
	}
}
