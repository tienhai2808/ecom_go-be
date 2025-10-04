package util

import (
	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/tienhai2808/ecom_go/internal/response"
)

func JSON(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, response.ApiResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func GenerateSlug(str string) string {
	return slug.Make(str)
}