package router

import (
	"backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewProductRouter(rg *gin.RouterGroup, productHandler handler.ProductHandler) {
	product := rg.Group(("/products")) 
	{
		product.GET("/all", productHandler.GetAllProducts)
	}
}