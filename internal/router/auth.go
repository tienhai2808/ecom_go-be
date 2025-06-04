package router

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(rg *gin.RouterGroup, config *config.AppConfig,authHandler handler.AuthHandler) {
	auth := rg.Group(("/auth"))
	{
		auth.POST("/signup", authHandler.Signup)
		auth.POST("/signup/verify", authHandler.VerifySignup)
		auth.POST("/signin", authHandler.Signin)
		auth.POST("/signout", middleware.RequireAuth(config), authHandler.Signout)
		auth.GET("/me", middleware.RequireAuth(config), authHandler.GetMe)
		auth.GET("/refresh-token", authHandler.RefreshToken)
		auth.POST("/forgot-password", authHandler.ForgotPassword)
	}
}