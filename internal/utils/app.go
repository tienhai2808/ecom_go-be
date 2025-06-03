package utils

import (
	"backend/internal/dto"

	"github.com/gin-gonic/gin"
)

func JSON(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, dto.ApiResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}
