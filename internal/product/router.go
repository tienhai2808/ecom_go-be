package product

import (
	"backend/internal/common"

	"github.com/gin-gonic/gin"
)

func ProductRouter(r *gin.RouterGroup, ctx *common.AppContext) {
	productGroup := r.Group("/products") 
	{
		productGroup.GET("/all", func(c *gin.Context) {
			c.String(200, "Tất cả sản phẩm đây")
		})
	}
}