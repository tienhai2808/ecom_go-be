package admin

// import (
// 	"backend/internal/auth"
// 	"backend/internal/common"

// 	"github.com/gin-gonic/gin"
// )

// func AdminRouter(r *gin.RouterGroup, ctx *common.AppContext) {
// 	adminService := NewService(ctx)
// 	handler := NewHandler(adminService, ctx)
// 	adminGroup := r.Group("/admin", auth.RequireAuth(ctx), RequireMultiRoles([]string{"admin"}))
// 	{
// 		adminGroup.GET("/users", handler.GetAllUsers)
// 		adminGroup.POST("/users", handler.CreateUser)
// 		adminGroup.PATCH("/users/:user_id", handler.UpdateUser)
// 		adminGroup.DELETE("/users", handler.DeleteUsers)
// 	} 
// }