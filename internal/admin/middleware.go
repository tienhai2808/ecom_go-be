package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"error":      "không có quyền truy cập",
			})
			return
		}

		role, _ := roleVal.(string)
		if role != "admin" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"error":      "không có quyền truy cập",
			})
			return
		}

		c.Next()
	}
}
