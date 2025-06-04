package middleware

import (
	"backend/internal/config"
	"backend/internal/dto"
	"backend/internal/repository"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(config *config.AppConfig, userRepository repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: 401,
				Message:    "Không có token",
			})
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("phương thức ký không hợp lệ: %v", token.Header["alg"])
			}
			return []byte(config.App.JWTAccessSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: 401,
				Message:    "Token không hợp lệ hoặc đã hết hạn",
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: 401,
				Message:    "Không thể lấy claims từ token",
			})
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: 401,
				Message:    "Không thể lấy user_id từ token",
			})
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: 401,
				Message:    "Không thể xác định quyền từ token",
			})
			return
		}

		ctx := c.Request.Context()
		user, err := userRepository.GetUserByID(ctx, userID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ApiResponse{
				StatusCode: 500,
				Message:    "Không thể xác lấy dữ liệu người dùng",
			})
			return
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: 401,
				Message:    "Không có quyền truy cập",
			})
			return
		}

		if string(user.Role) != role {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ApiResponse{
				StatusCode: 401,
				Message:    "Không có quyền truy cập",
			})
			return
		}

		c.Set("user", user)

		c.Next()
	}
}
