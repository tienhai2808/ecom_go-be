package router

import (
	"backend/internal/config"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"

	"github.com/gin-gonic/gin"
)

func NewProductRouter(rg *gin.RouterGroup, config *config.AppConfig, userRepository repository.UserRepository, productHandler handler.ProductHandler) {
	product := rg.Group("/products") 
	{
		product.GET("/all", productHandler.GetAllProducts)

		product.GET("/:product_id", productHandler.GetProductByID)

		product.POST("", middleware.RequireAuth(config, userRepository), middleware.RequireMultiRoles([]string{"admin"}), productHandler.CreateProduct)
	}
}