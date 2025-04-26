package admin

import (
	"backend/internal/auth"
	"backend/internal/common"

	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.RouterGroup, ctx *common.AppContext) {
	adminGroup := r.Group("/admin", auth.RequireAuth(ctx), RequireAdmin())
	{
		adminGroup.GET("/users", )
	} 
}