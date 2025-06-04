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
// 		authGroup.POST("/address", RequireAuth(ctx), handler.AddAddress)
// 		authGroup.PATCH("/address/:address_id", RequireAuth(ctx), handler.UpdateAddress)
// 		authGroup.DELETE("/address/:address_id", RequireAuth(ctx), handler.DeleteAddress)
// 	}
// }
