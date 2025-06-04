package middleware

import (
	"backend/internal/dto"
	"backend/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireMultiRoles(allowedRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: 401,
				Message:    "Không có quyền truy cập",
			})
			return
		}

		user, ok := userAny.(*model.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ApiResponse{
				StatusCode: 500,
				Message:    "Không thể chuyển đổi thông tin người dùng",
			})
			return
		}

		for _, allowedRole := range allowedRoles {
			if string(user.Role) == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
			StatusCode: 401,
			Message:    "Không có quyền truy cập",
		})
	}
}
