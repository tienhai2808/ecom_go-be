package auth

import (
	"e-commerce/internal/common"

	"github.com/gin-gonic/gin"
)

func AuthRouter(r *gin.RouterGroup, ctx *common.AppContext) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", SignupHandler(ctx))
		authGroup.POST("/verify-signup", VerifySignupHandler(ctx))
		authGroup.GET("/me", RequireAuth(ctx), GetMeHandler(ctx))
		authGroup.GET("/refresh-token", RefreshTokenHandler(ctx))
	}
}
