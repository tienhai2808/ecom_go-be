package handler

import (
	"backend/internal/common"
	"backend/internal/errors"
	"backend/internal/request"
	"backend/internal/service"
	"backend/internal/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Signup(c *gin.Context) {
	var req request.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	token, err := h.authService.Signup(req)
	if err != nil {
		switch err {
			case errors.ErrUsernameExists, errors.ErrEmailExists:
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
	var req request.VerifySignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	newUser, accessToken, refreshToken, err := h.authService.VerifySignup(req)
	if err != nil {
		switch err {
			case errors.ErrInvalidOTP, errors.ErrTooManyAttempts, errors.ErrEmailExists, errors.ErrUsernameExists:
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
	var req request.SigninRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	user, accessToken, refreshToken, err := h.authService.Signin(req)
	if err != nil {
		switch err {
			case errors.ErrIncorrectPassword, errors.ErrUserNotFound:
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
	userIDVal, exists := c.Get("user_id")
	if !exists {
		utils.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	userID, ok := userIDVal.(string)
	if !ok {
		utils.JSON(c, http.StatusUnauthorized, "Thông tin người dùng không hợp lệ", nil)
		return
	}

	user, err := h.authService.GetMe(userID)
	if err != nil {
		switch err {
			case errors.ErrUserNotFound:
				utils.JSON(c, http.StatusNotFound, err.Error(), nil)
			default:
				fmt.Printf("Lỗi ở GetMeService: %v\n", err)
				utils.JSON(c, http.StatusInternalServerError, "Không thể lấy thông tin người dùng", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Lấy thông tin người dùng thành công", gin.H{
		"user": user,
	})
}