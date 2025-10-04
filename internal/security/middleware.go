package security

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/mapper"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/response"
	"github.com/tienhai2808/ecom_go/internal/types"
)

func RequireAuth(accessName, secretKey string, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(accessName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    customErr.ErrUnAuth.Error(),
			})
			return
		}

		claims, err := ParseToken(tokenStr, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			})
			return
		}

		userID, userRole, err := ExtractToken(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			})
			return
		}

		user, err := userRepo.FindByIDWithProfile(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.ApiResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
			})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    customErr.ErrUserNotFound.Error(),
			})
			return
		}

		if user.Role != userRole {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    customErr.ErrInvalidUser.Error(),
			})
			return
		}

		userData := mapper.ToUserData(user)

		c.Set("user", userData)
		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAny, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    customErr.ErrUnAuth.Error(),
			})
			return
		}

		user, ok := userAny.(*types.UserData)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.ApiResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    "không thể chuyển đổi thông tin người dùng",
			})
			return
		}

		if user.Role != common.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, response.ApiResponse{
				StatusCode: http.StatusForbidden,
				Message:    customErr.ErrForbidden.Error(),
			})
			return
		}

		c.Next()
	}
}

func RequireRefreshToken(refreshName, secretKey string, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(refreshName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    customErr.ErrUnAuth.Error(),
			})
			return
		}

		claims, err := ParseToken(tokenStr, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			})
			return
		}

		userID, userRole, err := ExtractToken(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    err.Error(),
			})
			return
		}

		user, err := userRepo.FindByIDWithProfile(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, response.ApiResponse{
				StatusCode: http.StatusInternalServerError,
				Message:    err.Error(),
			})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    customErr.ErrUserNotFound.Error(),
			})
			return
		}

		if user.Role != userRole {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message:    customErr.ErrInvalidUser.Error(),
			})
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)
		c.Next()
	}
}
