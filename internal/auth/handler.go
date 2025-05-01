package auth

import (
	"backend/internal/common"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	service Service
	ctx     *common.AppContext
}

func NewHandler(service Service, ctx *common.AppContext) *Handler {
	return &Handler{
		service: service,
		ctx:     ctx,
	}
}

func (h *Handler) Test(c *gin.Context) {
	var req TestRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     err.Error(),
			"statusCode": http.StatusBadRequest,
		})
		return
	}
	fmt.Printf(`Nội dung request gửi đến: %s`, req.Request)
	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Hello world",
	})
}

func (h *Handler) Signup(c *gin.Context) {
	var req SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	token, err := h.service.Signup(req)
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

	c.JSON(http.StatusCreated, gin.H{
		"statusCode":        http.StatusCreated,
		"registrationToken": token,
		"message":           "Vui lòng kiểm tra email để lấy mã OTP",
	})
}

func (h *Handler) VerifySignup(c *gin.Context) {
	var req VerifySignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	user, accessToken, refreshToken, err := h.service.VerifySignup(req)
	if err != nil {
		switch err {
		case ErrUsernameExists, ErrEmailExists, ErrInvalidOTP, ErrTooManyAttempts, ErrTokenExpired:
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
		"user":       user,
	})
}

func (h *Handler) Signin(c *gin.Context) {
	var req SigninRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	user, accessToken, refreshToken, err := h.service.Signin(req)
	if err != nil {
		switch err {
		case ErrUserNotFound, ErrIncorrectPassword:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		default:
			fmt.Printf("Lỗi ở SigninService: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error":      "Không thể đăng nhập",
			})
		}
		return
	}

	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/ecom-go/auth/refresh-token", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"user":       user,
	})
}

func (h *Handler) Signout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/ecom-go/auth/refresh-token", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Đăng xuất thành công",
	})
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	token, err := h.service.ForgotPassword(req)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		default:
			fmt.Printf("Lỗi ở ForgotPasswordService: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error":      "Không thể đăng ký tài khoản",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode":          http.StatusOK,
		"forgotPasswordToken": token,
		"message":             "Vui lòng kiểm tra email để lấy mã OTP",
	})
}

func (h *Handler) VerifyForgotPassword(c *gin.Context) {
	var req VerifyForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	token, err := h.service.VerifyForgotPassword(req)
	if err != nil {
		switch err {
		case ErrInvalidOTP, ErrTooManyAttempts, ErrTokenExpired:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		default:
			fmt.Printf("Lỗi ở VerifyForgotPasswordService: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error":      "Không thể xác thực quên mật khẩu",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode":         http.StatusOK,
		"resetPasswordToken": token,
	})
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	user, accessToken, refreshToken, err := h.service.ResetPassword(req)
	if err != nil {
		switch err {
		case ErrUserNotFound, ErrTokenExpired:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		default:
			fmt.Printf("Lỗi ở ResetPasswordService: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error":      "Không thể đăng nhập",
			})
		}
		return
	}

	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/ecom-go/auth/refresh-token", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"user":       user,
	})
}

func (h *Handler) GetMe(c *gin.Context) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"error":      "không có quyền truy cập",
		})
		return
	}

	userID, _ := userIDVal.(string)

	user, err := h.service.GetMe(userID)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		default:
			fmt.Printf("Lỗi ở GetMeService: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error":      "Không thể lấy thông tin người dùng",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"user":       user,
	})
}

func (h *Handler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"error":      "Không có refresh token",
		})
		return
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.ctx.Config.App.JWTRefreshSecret), nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"error":      "Refresh token không hợp lệ hoặc đã hết hạn",
		})
		return
	}

	userID, ok1 := claims["user_id"].(string)
	role, ok2 := claims["role"].(string)
	if !ok1 || !ok2 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Claims không hợp lệ",
		})
		return
	}

	newAccessToken, err := GenerateToken(userID, role, 15*time.Minute, h.ctx.Config.App.JWTAccessSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error":      "Không thể tạo access token mới",
		})
		return
	}

	newRefreshToken, err := GenerateToken(userID, role, 24*7*time.Hour, h.ctx.Config.App.JWTRefreshSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error":      "Không thể tạo refresh token mới",
		})
	}

	c.SetCookie("access_token", newAccessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", newRefreshToken, 604800, "/ecom-go/auth/refresh-token", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"message":    "Token đã được làm mới",
	})
}

func (h *Handler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	currentUserIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"error":      "không có quyền truy cập",
		})
		return
	}
	currentUserID, _ := currentUserIDVal.(string)

	userID := c.Param("user_id")
	if currentUserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"error":      "không có quyền truy cập",
		})
		return
	}

	user, accessToken, refreshToken, err := h.service.ChangePassword(userID, req)
	if err != nil {
		switch err {
		case ErrUserNotFound, ErrIncorrectPassword:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		default:
			fmt.Printf("Lỗi ở ChangePasswordService: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error":      "Không thể lấy thông tin người dùng",
			})
		}
		return
	}

	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/ecom-go/auth/refresh-token", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"user":       user,
	})
}

func (h *Handler) UpdateUserInfo(c *gin.Context) {
	var req UpdateInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}
	userID := c.Param("user_id")
	currentUserIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"error":      "không có quyền truy cập",
		})
		return
	}

	currentUserID, _ := currentUserIDVal.(string)
	if currentUserID != userID {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"error":      "không có quyền truy cập",
		})
		return
	}

	updatedUser, err := h.service.UpdateInfo(userID, &req)
	if err != nil {
		switch err {
		case ErrUserNotFound, ErrUpdateFailed:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		default:
			fmt.Printf("Lỗi ở ChangePasswordService: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error":      "Không thể cập nhật thông tin người dùng",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"user":       updatedUser,
	})

}

func (h *Handler) AddAddress(c *gin.Context) {
	var req AddAddressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"error":      "không có quyền truy cập",
		})
		return
	}

	userID, _ := userIDVal.(string)
	address, err := h.service.AddAddress(userID, req)
	if err != nil {
		switch err {
		case ErrUserNotFound:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		default:
			fmt.Printf("Lỗi ở AddAddressService: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error":      "Không thể thêm địa chỉ",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"address":    address,
	})
}

func (h *Handler) UpdateAddress(c *gin.Context) {
	var req UpdateAddressRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": http.StatusBadRequest,
		})
		return
	}

	addressID := c.Param("address_id")

	userIDVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": http.StatusUnauthorized,
			"error":      "không có quyền truy cập",
		})
		return
	}

	userID, _ := userIDVal.(string)

	updatedAddress, err := h.service.UpdateAddress(userID, addressID, req)
	if err != nil {
		switch err {
		case ErrUnAuth, ErrAddressNotFound, ErrUpdateFailed:
			c.JSON(http.StatusBadRequest, gin.H{
				"statusCode": http.StatusBadRequest,
				"error":      err.Error(),
			})
		default:
			fmt.Printf("Lỗi ở UpdateAddressService: %v\n", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusInternalServerError,
				"error":      "Không thể cập nhật địa chỉ",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"address":    updatedAddress,
	})
}
