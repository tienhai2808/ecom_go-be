package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/mapper"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/types"
)

type ProfileHandler struct {
	profileSvc service.ProfileService
}

func NewProfileHandler(profileSvc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileSvc}
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	profileIDStr := c.Param("id")
	profileID, err := strconv.ParseInt(profileIDStr, 10, 64)
	if err != nil {
		common.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		common.JSON(c, http.StatusBadRequest, translated, nil)
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		common.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*types.UserData)
	if !ok {
		common.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	if user.Profile.ID != profileID {
		common.JSON(c, http.StatusForbidden, "Không có quyền truy cập", nil)
		return
	}

	userRes, err := h.profileSvc.UpdateProfile(ctx, user.Profile.ID, &req)
	if err != nil {
		switch err {
		case customErr.ErrProfileNotFound, customErr.ErrUserNotFound:
			common.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			common.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	common.JSON(c, http.StatusOK, "Cập nhật hồ sơ người dùng thành công", gin.H{
		"user": mapper.ToUserResponse(userRes),
	})
}
