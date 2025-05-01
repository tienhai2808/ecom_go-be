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
		authGroup.POST("/test", handler.Test);
		authGroup.POST("/signup", handler.Signup)
		authGroup.POST("/signup/verify-email", handler.VerifySignup)
		authGroup.POST("/signin", handler.Signin)
		authGroup.POST("/signout", RequireAuth(ctx), handler.Signout)
		authGroup.POST("/forgot-password", handler.ForgotPassword)
		authGroup.POST("/forgot-password/verify-email", handler.VerifyForgotPassword)
		authGroup.POST("/reset-password", handler.ResetPassword)
		authGroup.GET("/me", RequireAuth(ctx), handler.GetMe)
		authGroup.PATCH("/update-info/:user_id", RequireAuth(ctx), handler.UpdateUserInfo)
		authGroup.PUT("/change-password/:user_id", RequireAuth(ctx), handler.ChangePassword)
		authGroup.GET("/refresh-token", handler.RefreshToken)
		authGroup.POST("/add-address", RequireAuth(ctx), handler.AddAddress)
	}
}
