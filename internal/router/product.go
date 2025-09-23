package router

import (
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/common"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"

	"github.com/gin-gonic/gin"
)

func NewProductRouter(rg *gin.RouterGroup, config *config.Config, userRepository repository.UserRepository, productHandler *handler.ProductHandler) {
	accessName := config.App.AccessName
	secretKey := config.App.JWTSecret

	product := rg.Group("/products")
	{
		product.GET("/all", productHandler.GetAllProducts)

		product.GET("/:product_id", productHandler.GetProductByID)

		product.POST("", security.RequireAuthAndRole(accessName, secretKey, common.RoleAdmin, userRepository), productHandler.CreateProduct)

		product.PATCH("/:product_id", security.RequireAuthAndRole(accessName, secretKey, common.RoleAdmin, userRepository), productHandler.UpdateProduct)

		product.DELETE("/:product_id", security.RequireAuthAndRole(accessName, secretKey, common.RoleAdmin, userRepository), productHandler.DeleteProduct)

		product.DELETE("/many", security.RequireAuthAndRole(accessName, secretKey, common.RoleAdmin, userRepository), productHandler.DeleteManyProducts)
	}
}
