package handler

import (
	"backend/internal/common"
	"backend/internal/config"
	customErr "backend/internal/errors"
	"backend/internal/model"
	"backend/internal/request"
	"backend/internal/service"
	"backend/internal/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthHandler struct {
	authService service.AuthService
	config      *config.AppConfig
}

func NewAuthHandler(authService service.AuthService, config *config.AppConfig) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      config,
	}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	token, err := h.authService.Signup(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUsernameExists, customErr.ErrEmailExists:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở SignupService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể đăng ký tài khoản", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Vui lòng kiểm tra email để lấy mã OTP", gin.H{
		"registration_token": token,
	})
}

func (h *AuthHandler) VerifySignup(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.VerifySignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	newUser, accessToken, refreshToken, err := h.authService.VerifySignup(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrInvalidOTP, customErr.ErrTooManyAttempts, customErr.ErrEmailExists, customErr.ErrUsernameExists, customErr.ErrKeyNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở VerifySignupService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể xác minh đăng ký", nil)
		}
		return
	}

	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/ecom-go/auth/refresh-token", "", false, true)

	utils.JSON(c, http.StatusOK, "Đăng ký thành công", gin.H{
		"user": newUser,
	})
}

func (h *AuthHandler) Signin(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.SigninRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	user, accessToken, refreshToken, err := h.authService.Signin(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrIncorrectPassword, customErr.ErrUserNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở SigninService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể đăng nhập", nil)
		}
		return
	}

	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/ecom-go/auth/refresh-token", "", false, true)

	utils.JSON(c, http.StatusOK, "Đăng nhập thành công", gin.H{
		"user": user,
	})
}

func (h *AuthHandler) Signout(c *gin.Context) {
	c.SetCookie("access_token", "", -1, "/", "", false, true)
	c.SetCookie("refresh_token", "", -1, "/ecom-go/auth/refresh-token", "", false, true)

	utils.JSON(c, http.StatusOK, "Đăng xuất thành công", nil)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userAny, exists := c.Get("user")
	if !exists {
		utils.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		utils.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Lấy thông tin người dùng thành công", gin.H{
		"user": user,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	ctx := c.Request.Context()

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		utils.JSON(c, http.StatusUnauthorized, "Không có refresh token", nil)
		return
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(h.config.App.JWTRefreshSecret), nil
	})
	if err != nil {
		utils.JSON(c, http.StatusUnauthorized, "Refresh token không hợp lệ hoặc đã hết hạn", nil)
		return
	}

	userID, ok1 := claims["user_id"].(string)
	role, ok2 := claims["role"].(string)
	if !ok1 || !ok2 {
		utils.JSON(c, http.StatusUnauthorized, "Không thể lấy thông tin người dùng từ refresh token", nil)
		return
	}

	_, err = h.authService.GetMe(ctx, userID)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			utils.JSON(c, http.StatusUnauthorized, "Yêu cầu không hợp lệ", nil)
		default:
			fmt.Printf("Lỗi ở GetMeService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể làm mới token", nil)
		}
		return
	}

	newAccessToken, err := utils.GenerateToken(userID, role, 15*time.Minute, h.config.App.JWTAccessSecret)
	if err != nil {
		fmt.Printf("Lỗi khi tạo access token mới: %v\n", err)
		utils.JSON(c, http.StatusInternalServerError, "Không thể tạo access token mới", nil)
		return
	}

	newRefreshToken, err := utils.GenerateToken(userID, role, 7*24*time.Hour, h.config.App.JWTRefreshSecret)
	if err != nil {
		fmt.Printf("Lỗi khi tạo refresh token mới: %v\n", err)
		utils.JSON(c, http.StatusInternalServerError, "Không thể tạo refresh token mới", nil)
		return
	}

	c.SetCookie("access_token", newAccessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", newRefreshToken, 604800, "/ecom-go/auth/refresh-token", "", false, true)

	utils.JSON(c, http.StatusOK, "Làm mới token thành công", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	token, err := h.authService.ForgotPassword(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở ForgotPasswordService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể lấy lại mật khẩu", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Vui lòng kiểm tra email để lấy mã OTP", gin.H{
		"forgot_password_token": token,
	})
}

func (h *AuthHandler) VerifyForgotPassword(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.VerifyForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	token, err := h.authService.VerifyForgotPassword(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrInvalidOTP, customErr.ErrKeyNotFound, customErr.ErrTooManyAttempts:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở VerifyForgotPasswordService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể xác thực quên mật khẩu", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Xác thực quên mật khẩu thành công", gin.H{
		"reset_password_token": token,
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.ResetPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	user, accessToken, refreshToken, err := h.authService.ResetPassword(ctx, req)
	if err != nil {
		fmt.Println(err)
		switch err {
		case customErr.ErrUserNotFound, customErr.ErrKeyNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở ResetPasswordService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể làm mới mật khẩu", nil)
		}
		return
	}

	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/ecom-go/auth/refresh-token", "", false, true)

	utils.JSON(c, http.StatusOK, "Lấy lại mật khẩu thành công", gin.H{
		"user": user,
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		utils.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		utils.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	user, accessToken, refreshToken, err := h.authService.ChangePassword(ctx, user, req)
	if err != nil {
		switch err {
		case customErr.ErrIncorrectPassword, customErr.ErrUserNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở ChangePasswordService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể thay đổi mật khẩu", nil)
		}
		return
	}

	c.SetCookie("access_token", accessToken, 900, "/", "", false, true)
	c.SetCookie("refresh_token", refreshToken, 604800, "/ecom-go/auth/refresh-token", "", false, true)

	utils.JSON(c, http.StatusOK, "Thay đổi mật khẩu thành công", gin.H{
		"user": user,
	})
}

func (h *AuthHandler) UpdateUserProfile(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		utils.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		utils.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	reqUserID := c.Param("user_id")
	if user.ID != reqUserID {
		utils.JSON(c, http.StatusUnauthorized, "Không có quyền truy cập", nil)
		return
	}

	updatedUser, err := h.authService.UpdateUserProfile(ctx, user, &req)
	if err != nil {
		switch err {
		case customErr.ErrUserProfileNotFound, customErr.ErrUserNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở UpdateUserProfileService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể thay đổi hồ sơ người dùng", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Cập nhật hồ sơ người dùng thành công", gin.H{
		"user": updatedUser,
	})
}
