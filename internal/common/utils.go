package common

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
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

func IsUniqueViolation(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062 
	}
	return false
}
