package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/security"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/utils"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
	userService service.UserService
	config      *config.Config
}

func NewAuthHandler(authService service.AuthService, userService service.UserService, config *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		config:      config,
	}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, translated, nil)
		return
	}

	token, err := h.authService.Signup(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUsernameExists, customErr.ErrEmailExists:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
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
		utils.JSON(c, http.StatusBadRequest, translated, nil)
		return
	}

	userRes, accessToken, refreshToken, err := h.authService.VerifySignup(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrInvalidOTP, customErr.ErrTooManyAttempts, customErr.ErrEmailExists, customErr.ErrUsernameExists, customErr.ErrKeyNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	c.SetCookie(h.config.App.AccessName, accessToken, 900, "/", "", false, true)
	c.SetCookie(h.config.App.RefreshName, refreshToken, 604800, fmt.Sprintf("%s/auth/refresh-token", h.config.App.ApiPrefix), "", false, true)

	utils.JSON(c, http.StatusOK, "Đăng ký thành công", gin.H{
		"user": userRes,
	})
}

func (h *AuthHandler) Signin(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.SigninRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, translated, nil)
		return
	}

	userRes, accessToken, refreshToken, err := h.authService.Signin(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrIncorrectPassword, customErr.ErrUserNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	c.SetCookie(h.config.App.AccessName, accessToken, 900, "/", "", false, true)
	c.SetCookie(h.config.App.RefreshName, refreshToken, 604800, fmt.Sprintf("%s/auth/refresh-token", h.config.App.ApiPrefix), "", false, true)

	utils.JSON(c, http.StatusOK, "Đăng nhập thành công", gin.H{
		"user": userRes,
	})
}

func (h *AuthHandler) Signout(c *gin.Context) {
	c.SetCookie(h.config.App.AccessName, "", -1, "/", "", false, true)
	c.SetCookie(h.config.App.RefreshName, "", -1, fmt.Sprintf("%s/auth/refresh-token", h.config.App.ApiPrefix), "", false, true)

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

	userRes := h.authService.ConvertToDto(user)

	utils.JSON(c, http.StatusOK, "Lấy thông tin người dùng thành công", gin.H{
		"user": userRes,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		utils.JSON(c, http.StatusUnauthorized, customErr.ErrUserIdNotFound.Error(), nil)
		return
	}

	userRole := c.GetString("user_role")
	if userRole == "" {
		utils.JSON(c, http.StatusUnauthorized, customErr.ErrUserRoleNotFound.Error(), nil)
		return
	}

	newAccessToken, err := security.GenerateToken(userID, userRole, 15*time.Minute, h.config.App.JWTSecret)
	if err != nil {
		utils.JSON(c, http.StatusInternalServerError, fmt.Sprintf("tạo access token mới thất bại: %v", err), nil)
		return
	}

	newRefreshToken, err := security.GenerateToken(userID, userRole, 7*24*time.Hour, h.config.App.JWTSecret)
	if err != nil {
		utils.JSON(c, http.StatusInternalServerError, fmt.Sprintf("tạo refresh token mới thất bại: %v", err), nil)
		return
	}

	c.SetCookie(h.config.App.AccessName, newAccessToken, 900, "/", "", false, true)
	c.SetCookie(h.config.App.RefreshName, newRefreshToken, 604800, fmt.Sprintf("%s/auth/refresh-token", h.config.App.ApiPrefix), "", false, true)

	utils.JSON(c, http.StatusOK, "Làm mới token thành công", nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.ForgotPasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, translated, nil)
		return
	}

	token, err := h.authService.ForgotPassword(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			utils.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
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
		utils.JSON(c, http.StatusBadRequest, translated, nil)
		return
	}

	token, err := h.authService.VerifyForgotPassword(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrInvalidOTP, customErr.ErrKeyNotFound, customErr.ErrTooManyAttempts:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
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
		utils.JSON(c, http.StatusBadRequest, translated, nil)
		return
	}

	userRes, accessToken, refreshToken, err := h.authService.ResetPassword(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound, customErr.ErrKeyNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	c.SetCookie(h.config.App.AccessName, accessToken, 900, "/", "", false, true)
	c.SetCookie(h.config.App.RefreshName, refreshToken, 604800, fmt.Sprintf("%s/auth/refresh-token", h.config.App.ApiPrefix), "", false, true)

	utils.JSON(c, http.StatusOK, "Lấy lại mật khẩu thành công", gin.H{
		"user": userRes,
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.ChangePasswordRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, translated, nil)
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

	userRes, accessToken, refreshToken, err := h.authService.ChangePassword(ctx, user, req)
	if err != nil {
		switch err {
		case customErr.ErrIncorrectPassword, customErr.ErrUserNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	c.SetCookie(h.config.App.AccessName, accessToken, 900, "/", "", false, true)
	c.SetCookie(h.config.App.RefreshName, refreshToken, 604800, fmt.Sprintf("%s/auth/refresh-token", h.config.App.ApiPrefix), "", false, true)

	utils.JSON(c, http.StatusOK, "Thay đổi mật khẩu thành công", gin.H{
		"user": userRes,
	})
}

func (h *AuthHandler) UpdateUserProfile(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.UpdateProfileRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, translated, nil)
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

	userRes, err := h.authService.UpdateUserProfile(ctx, user, &req)
	if err != nil {
		switch err {
		case customErr.ErrUserProfileNotFound, customErr.ErrUserNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			utils.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Cập nhật hồ sơ người dùng thành công", gin.H{
		"user": userRes,
	})
}
