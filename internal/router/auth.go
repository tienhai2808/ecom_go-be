package router

import (
	"github.com/tienhai2808/ecom_go/config"
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
		auth.POST("/signup", authHandler.SignUp)

		auth.POST("/signup/verify", authHandler.VerifySignUp)

		auth.POST("/signin", authHandler.SignIn)

		auth.POST("/signout", security.RequireAuth(accessName, secretKey, userRepository), authHandler.SignOut)

		auth.GET("/me", security.RequireAuth(accessName, secretKey, userRepository), authHandler.GetMe)

		auth.GET("/refresh-token", security.RequireRefreshToken(refreshName, secretKey, userRepository), authHandler.RefreshToken)

		auth.POST("/forgot-password", authHandler.ForgotPassword)

		auth.POST("/forgot-password/verify", authHandler.VerifyForgotPassword)

		auth.POST("/reset-password", authHandler.ResetPassword)

		auth.POST("/change-password", security.RequireAuth(accessName, secretKey, userRepository), authHandler.ChangePassword)
	}
}
