package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireMultiRoles(allowedRoles []string) gin.HandlerFunc {
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

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"error":      "không có quyền truy cập",
		})
	}
}
