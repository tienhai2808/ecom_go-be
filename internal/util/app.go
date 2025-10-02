package util

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gosimple/slug"
	"github.com/tienhai2808/ecom_go/internal/dto"
	"github.com/tienhai2808/ecom_go/internal/initialization"
)

func JSON(c *gin.Context, statusCode int, message string, data any) {
	c.JSON(statusCode, dto.ApiResponse{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func NewSnowflakeID() (int64, error) {
	id, err := initialization.Sf.NextID()
	if err != nil {
		return 0, fmt.Errorf("tạo ID thất bại")
	}

	return id, nil
}

func GenerateSlug(str string) string {
	return slug.Make(str)
}