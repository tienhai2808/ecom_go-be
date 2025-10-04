package router

import (
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(rg *gin.RouterGroup, cfg *config.Config, userRepo repository.UserRepository, authHdl *handler.AuthHandler) {
	accessName := cfg.App.AccessName
	refreshName := cfg.App.RefreshName
	secretKey := cfg.App.JWTSecret

	auth := rg.Group("/auth")
	{
		auth.POST("/signup", authHdl.SignUp)

		auth.POST("/signup/verify", authHdl.VerifySignUp)

		auth.POST("/signin", authHdl.SignIn)

		auth.POST("/signout", security.RequireAuth(accessName, secretKey, userRepo), authHdl.SignOut)

		auth.GET("/me", security.RequireAuth(accessName, secretKey, userRepo), authHdl.GetMe)

		auth.GET("/refresh-token", security.RequireRefreshToken(refreshName, secretKey, userRepo), authHdl.RefreshToken)

		auth.POST("/forgot-password", authHdl.ForgotPassword)

		auth.POST("/forgot-password/verify", authHdl.VerifyForgotPassword)

		auth.POST("/reset-password", authHdl.ResetPassword)

		auth.POST("/change-password", security.RequireAuth(accessName, secretKey, userRepo), authHdl.ChangePassword)
	}
}
