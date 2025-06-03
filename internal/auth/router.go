package auth

// import (
// 	"backend/internal/common"

// 	"github.com/gin-gonic/gin"
// )

// func AuthRouter(r *gin.RouterGroup, ctx *common.AppContext) {
// 	authService := NewService(ctx)
// 	handler := NewHandler(authService, ctx)

// 	authGroup := r.Group("/auth")
// 	{
// 		authGroup.POST("/test", handler.Test);
// 		authGroup.POST("/forgot-password", handler.ForgotPassword)
// 		authGroup.POST("/forgot-password/verify-email", handler.VerifyForgotPassword)
// 		authGroup.POST("/reset-password", handler.ResetPassword)
// 		authGroup.PATCH("/update-info/:user_id", RequireAuth(ctx), handler.UpdateUserInfo)
// 		authGroup.PUT("/change-password/:user_id", RequireAuth(ctx), handler.ChangePassword)
// 		authGroup.GET("/refresh-token", handler.RefreshToken)
// 		authGroup.GET("/addresses", RequireAuth(ctx), handler.GetAddresses)
// 		authGroup.POST("/address", RequireAuth(ctx), handler.AddAddress)
// 		authGroup.PATCH("/address/:address_id", RequireAuth(ctx), handler.UpdateAddress)
// 		authGroup.DELETE("/address/:address_id", RequireAuth(ctx), handler.DeleteAddress)
// 	}
// }
