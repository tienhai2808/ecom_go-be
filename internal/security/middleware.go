package security

import (
	"net/http"

	"github.com/tienhai2808/ecom_go/internal/dto"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/gin-gonic/gin"
)

func RequireAuthAndRole(accessName, secretKey, roleName string, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(accessName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message: customErr.ErrUnAuth.Error(),
			})
			return
		}

		claims, err := ParseToken(tokenStr, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message: err.Error(),
			})
			return
		}

		userID, userRole, err := ExtractToken(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message: err.Error(),
			})
			return
		}

		user, err := userRepo.FindByIDWithProfile(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ApiResponse{
				StatusCode: http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message: customErr.ErrUserNotFound.Error(),
			})
			return
		}

		if user.Role != userRole || user.Role != roleName {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message: customErr.ErrInvalidUser.Error(),
			})
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func RequireRefreshToken(refreshName, secretKey string, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(refreshName)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message: customErr.ErrUnAuth.Error(),
			})
			return
		}

		claims, err := ParseToken(tokenStr, secretKey)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message: err.Error(),
			})
			return
		}

		userID, userRole, err := ExtractToken(claims)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message: err.Error(),
			})
			return
		}

		user, err := userRepo.FindByIDWithProfile(c.Request.Context(), userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ApiResponse{
				StatusCode: http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}
		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message: customErr.ErrUserNotFound.Error(),
			})
			return
		}

		if user.Role != userRole {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: http.StatusUnauthorized,
				Message: customErr.ErrInvalidUser.Error(),
			})
			return
		}

		c.Set("user_id", user.ID)
		c.Set("user_role", user.Role)
		c.Next()
	}
}
