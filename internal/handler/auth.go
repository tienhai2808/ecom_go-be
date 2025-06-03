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
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", translated)
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

	utils.JSON(c, http.StatusOK, "Vui lòng kiểm tra email để lấy mã OTP", token)
}

func (h *AuthHandler) VerifySignup(c *gin.Context) {
	var req request.VerifySignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", translated)
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
	
	utils.JSON(c, http.StatusOK, "Đăng ký thành công", newUser)
}