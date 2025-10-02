package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/tienhai2808/ecom_go/internal/common"
	customErr "github.com/tienhai2808/ecom_go/internal/errors"
	"github.com/tienhai2808/ecom_go/internal/model"
	"github.com/tienhai2808/ecom_go/internal/request"
	"github.com/tienhai2808/ecom_go/internal/service"
	"github.com/tienhai2808/ecom_go/internal/util"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	util.JSON(c, http.StatusOK, "Lấy danh sách người dùng thành công", gin.H{
		"users": users,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		util.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			util.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	util.JSON(c, http.StatusOK, "Lấy người dùng thành công", gin.H{
		"user": user,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		util.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	newUser, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUsernameExists, customErr.ErrEmailExists:
			util.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	util.JSON(c, http.StatusCreated, "Thêm mới người dùng thành công", gin.H{
		"user": newUser,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		util.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	userIDStr := c.Param("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		util.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}

	updatedUser, err := h.userService.UpdateUser(ctx, userID, &req)
	if err != nil {
		switch err {
		case customErr.ErrEmailExists, customErr.ErrUsernameExists, customErr.ErrUserNotFound, customErr.ErrUserProfileNotFound:
			util.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	util.JSON(c, http.StatusOK, "Cập nhật người dùng thành công", gin.H{
		"user": updatedUser,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	userAny, exists := c.Get("user")
	if !exists {
		util.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		util.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	reqUserIDStr := c.Param("id")
	reqUserID, err := strconv.ParseInt(reqUserIDStr, 10, 64)
	if err != nil {
		util.JSON(c, http.StatusBadRequest, customErr.ErrInvalidID.Error(), nil)
		return
	}
	
	if reqUserID == user.ID {
		util.JSON(c, http.StatusConflict, "Không thể xóa chính bạn", nil)
		return
	}

	if err := h.userService.DeleteUser(ctx, reqUserID); err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			util.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	util.JSON(c, http.StatusOK, "Xóa người dùng thành công", nil)
}

func (h *UserHandler) DeleteManyUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()
	var req request.DeleteManyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		util.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	userAny, exists := c.Get("user")
	if !exists {
		util.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		util.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	rowsAccepted, err := h.userService.DeleteUsers(ctx, user.ID, req)
	if err != nil {
		switch err {
		case customErr.ErrUserConflict:
			util.JSON(c, http.StatusConflict, err.Error(), nil)
		default:
			util.JSON(c, http.StatusInternalServerError, err.Error(), nil)
		}
		return
	}

	message := fmt.Sprintf("Xóa thành công %d người dùng", rowsAccepted)
	util.JSON(c, http.StatusOK, message, nil)
}
