package router

import (
	"backend/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(rg *gin.RouterGroup, config *config.Config, userRepository repository.UserRepository, authHandler *handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/signup", authHandler.Signup)

		auth.POST("/signup/verify", authHandler.VerifySignup)

		auth.POST("/signin", authHandler.Signin)

		auth.POST("/signout", middleware.RequireAuth(config, userRepository), authHandler.Signout)

		auth.GET("/me", middleware.RequireAuth(config, userRepository), authHandler.GetMe)

		auth.GET("/refresh-token", authHandler.RefreshToken)

		auth.POST("/forgot-password", authHandler.ForgotPassword)

		auth.POST("/forgot-password/verify", authHandler.VerifyForgotPassword)

		auth.POST("/reset-password", authHandler.ResetPassword)

		auth.PUT("/change-password/:user_id", middleware.RequireAuth(config, userRepository), authHandler.ChangePassword)

		auth.PATCH("/update-info/:user_id", middleware.RequireAuth(config, userRepository), authHandler.UpdateUserProfile)
	}
}
