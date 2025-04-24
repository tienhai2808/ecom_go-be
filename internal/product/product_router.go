package product

import (
	"github.com/gin-gonic/gin"
)

func ProductRouter(r *gin.RouterGroup) {
	productGroup := r.Group("/products") 
	{
		productGroup.GET("/all", func(c *gin.Context) {
			c.String(200, "Tất cả sản phẩm đây")
		})
	}
}