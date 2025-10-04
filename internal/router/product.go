package router

import (
	"github.com/tienhai2808/ecom_go/config"
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
		product.GET("", productHandler.GetAllProducts)

		product.GET("/:id", productHandler.GetProductByID)

		product.POST("", security.RequireAuth(accessName, secretKey, userRepository), security.RequireAdmin(), productHandler.CreateProduct)

		product.PATCH("/:id", security.RequireAuth(accessName, secretKey, userRepository), security.RequireAdmin(), productHandler.UpdateProduct)

		product.DELETE("/:id", security.RequireAuth(accessName, secretKey, userRepository), security.RequireAdmin(), productHandler.DeleteProduct)

		product.DELETE("", security.RequireAuth(accessName, secretKey, userRepository), security.RequireAdmin(), productHandler.DeleteManyProducts)
	}
}
