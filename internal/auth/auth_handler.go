package auth

import (
	"e-commerce/internal/common"
	"e-commerce/internal/database"
	"e-commerce/internal/user"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SignupHandler(c *gin.Context) {
	var req SignupSchema

	if err := c.ShouldBindJSON(&req); err != nil {
		translated := common.HandleValidationError(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"errors":     translated,
			"statusCode": 400,
		})
		return
	}

	db, err := database.ConnectToDatabase()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error": "Không thể kết nối đến cơ sở dữ liệu",
		})
		return 
	}

	var existingUser user.User
	err = db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error
	if err == nil {
		if existingUser.Username == req.Username {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusBadRequest,
				"error": "Username đã tồn tại",
			})
			return
		}
		if existingUser.Email == req.Email {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": http.StatusBadRequest,
				"error": "Email đã tồn tại",
			})
			return
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error": "Không thể tạo người dùng",
		})
		return
	}

	user, err := SignupService(db, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": http.StatusInternalServerError,
			"error": "Không thể tạo người dùng",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": http.StatusOK,
		"user": gin.H{
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
