package auth

import (
	"e-commerce/internal/common"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func SignupHandler(ctx *common.AppContext) gin.HandlerFunc {
	return func (c *gin.Context)  {
		var req SignupSchema

		if err := c.ShouldBindJSON(&req); err != nil {
			translated := common.HandleValidationError(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"errors":     translated,
				"statusCode": 400,
			})
			return
		}
	
		token, err := SignupService(ctx, req)
		if err != nil {
			switch err {
			case ErrUsernameExists, ErrEmailExists:
				c.JSON(http.StatusBadRequest, gin.H{
					"statusCode": http.StatusBadRequest,
					"error":      err.Error(),
				})
			default:
				fmt.Printf("Lỗi ở SignupService: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"statusCode": http.StatusInternalServerError,
					"error":      "Không thể đăng ký tài khoản",
				})
			}
			return
		}
	
		c.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"registrationToken": token,
			"message": "Vui lòng kiểm tra email để lấy mã OTP",
		})
	}
}

func VerifySignupHandler(ctx *common.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req VerifySignupSchema

		if err := c.ShouldBindJSON(&req); err != nil {
			translated := common.HandleValidationError(err)
			c.JSON(http.StatusBadRequest, gin.H{
				"errors":     translated,
				"statusCode": 400,
			})
			return
		}

		user, accessToken, refreshToken, err := VerifySignupService(ctx, req)
		if err != nil {
			switch err {
			case ErrUsernameExists, ErrEmailExists:
				c.JSON(http.StatusBadRequest, gin.H{
					"statusCode": http.StatusBadRequest,
					"error":      err.Error(),
				})
			default:
				fmt.Printf("Lỗi ở VerifySignupService: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"statusCode": http.StatusInternalServerError,
					"error":      "Không thể đăng ký tài khoản",
				})
			}
			return
		}

		c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
		c.SetCookie("refresh_token", refreshToken, 604800, "/ecom-go/auth/refresh-token", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"user": user,
		})
	}
}

func GetMeHandler(ctx *common.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"error": "Không tìm thấy thông tin người dùng",
			})
			return
		}

		userID, ok := userIDVal.(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"error":      "ID người dùng không hợp lệ",
			})
			return
		}

		user, err := GetMeService(ctx, userID)
		if err != nil {
			if err.Error() == "người dùng không tồn tại" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"statusCode": http.StatusUnauthorized,
					"error":      err.Error(),
				})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"statusCode": http.StatusInternalServerError,
					"error":      "Không thể lấy thông tin người dùng",
				})
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"user": user,
		})
	}
}

func RefreshTokenHandler(ctx *common.AppContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken, err := c.Cookie("refresh_token")
		if err != nil || refreshToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"error": "Không có refresh_token",
			})
			return
		}

		claims := jwt.MapClaims{}
		_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(ctx.Config.App.JWTRefreshSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"statusCode": http.StatusUnauthorized,
				"error": "Refresh token không hợp lệ hoặc đã hết hạn",
			})
			return
		}

		userID, ok1 := claims["user_id"].(string)
		role, ok2 := claims["role"].(string)
		if !ok1 || !ok2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Claims không hợp lệ"})
			return
		}

		newAccessToken, err := GenerateToken(userID, role, 15*time.Minute, ctx.Config.App.JWTAccessSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error": "Không thể tạo access token mới",
			})
			return
		}

		c.SetCookie("access_token", newAccessToken, 900, "/", "", false, true)
		c.SetCookie("refresh_token", refreshToken, 604800, "/ecom-go/auth/refresh-token", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"statusCode": http.StatusOK,
			"message": "Access token mới đã được cấp phát",
		})
	}
}
