package auth

import "github.com/gin-gonic/gin"

func AuthRouter(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/signup", SignupHandler)
	}
}