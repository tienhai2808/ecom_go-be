package user

import (
	"e-commerce/internal/common"

	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup, ctx *common.AppContext) {
	userGroup := r.Group("/users") 
	{
		userGroup.GET("/all", func(c *gin.Context) {
			c.String(200, "Hello world!!!")
		})
	}
}