package auth

import (
	"backend/internal/common"

	"github.com/gin-gonic/gin"
)

func AuthRouter(r *gin.RouterGroup, ctx *common.AppContext) {
	authService := NewService(ctx)
	handler := NewHandler(authService, ctx)

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", handler.Signup)
		authGroup.POST("/verify-signup", handler.VerifySignup)
		authGroup.POST("/signin", handler.Signin)
		authGroup.POST("/signout", RequireAuth(ctx), handler.Signout)
		authGroup.POST("/forgot-password", handler.ForgotPassword)
		authGroup.POST("/verify-forgot-password", handler.VerifyForgotPassword)
		authGroup.POST("/reset-password", handler.ResetPassword)
		authGroup.GET("/me", RequireAuth(ctx), handler.GetMe)
		authGroup.GET("/refresh-token", handler.RefreshToken)
	}
} 
