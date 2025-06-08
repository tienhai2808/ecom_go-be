package handler

import (
	"backend/internal/common"
	customErr "backend/internal/errors"
	"backend/internal/model"
	"backend/internal/request"
	"backend/internal/service"
	"backend/internal/utils"
	"fmt"
	"net/http"

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
	ctx := c.Request.Context()

	users, err := h.userService.GetAllUsers(ctx)
	if err != nil {
		fmt.Printf("Lỗi ở GetAllUserService: %v\n", err)
		utils.JSON(c, http.StatusInternalServerError, "Không thế lấy danh sách người dùng", nil)
		return
	}

	utils.JSON(c, http.StatusOK, "Lấy danh sách người dùng thành công", gin.H{
		"users": users,
	})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("user_id")

	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			utils.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở GetUserByIDService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể lấy người dùng", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Lấy người dùng thành công", gin.H{
		"user": user,
	})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	newUser, err := h.userService.CreateUser(ctx, req)
	if err != nil {
		switch err {
		case customErr.ErrUsernameExists, customErr.ErrEmailExists:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở CreateUserService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể thêm mới người dùng", nil)
		}
		return
	}

	utils.JSON(c, http.StatusCreated, "Thêm mới người dùng thành công", gin.H{
		"user": newUser,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.UpdateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		utils.JSON(c, http.StatusBadRequest, "Dữ liệu gửi lên không hợp lệ", gin.H{
			"errors": translated,
		})
		return
	}

	userID := c.Param("user_id")

	updatedUser, err := h.userService.UpdateUser(ctx, userID, &req)
	if err != nil {
		switch err {
		case customErr.ErrEmailExists, customErr.ErrUsernameExists, customErr.ErrUserNotFound, customErr.ErrUserProfileNotFound:
			utils.JSON(c, http.StatusBadRequest, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở UpdateUserService %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể cập nhật người dùng", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Cập nhật người dùng thành công", gin.H{
		"user": updatedUser,
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	userAny, exists := c.Get("user")
	if !exists {
		utils.JSON(c, http.StatusUnauthorized, "Không có thông tin người dùng", nil)
		return
	}

	user, ok := userAny.(*model.User)
	if !ok {
		fmt.Println("Lỗi ở đổi kiểu dữ liệu user lấy từ context")
		utils.JSON(c, http.StatusInternalServerError, "Không thể chuyển đổi thông tin người dùng", nil)
		return
	}

	reqUserID := c.Param("user_id")
	if reqUserID == user.ID {
		utils.JSON(c, http.StatusConflict, "Không thể xóa chính bạn", nil)
		return
	}

	if err := h.userService.DeleteUserByID(ctx, reqUserID); err != nil {
		switch err {
		case customErr.ErrUserNotFound:
			utils.JSON(c, http.StatusNotFound, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở DeleteUserByIDService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể xóa người dùng", nil)
		}
		return
	}

	utils.JSON(c, http.StatusOK, "Xóa người dùng thành công", nil)
}

func (h *UserHandler) DeleteManyUsers(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.DeleteManyUsersRequest

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

	rowsAccepted, err := h.userService.DeleteManyUsers(ctx, user.ID, req)
	if err != nil {
		switch err {
		case customErr.ErrUserConflict:
			utils.JSON(c, http.StatusConflict, err.Error(), nil)
		default:
			fmt.Printf("Lỗi ở DeleteManyUsersService: %v\n", err)
			utils.JSON(c, http.StatusInternalServerError, "Không thể xóa người dùng", nil)
		}
		return
	}

	message := fmt.Sprintf("Xóa thành công %d người dùng", rowsAccepted)
	utils.JSON(c, http.StatusOK, message, nil)
}
