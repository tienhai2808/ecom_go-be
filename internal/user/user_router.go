package user

import (
	"github.com/gin-gonic/gin"
)

func UserRouter(r *gin.RouterGroup) {
	userGroup := r.Group("/users") 
	{
		userGroup.GET("/me", MeHandler)
		userGroup.GET("/all", func(c *gin.Context) {
			c.String(200, "Hello world!!!")
		})
	}
}