package router

import (
	"github.com/tienhai2808/ecom_go/internal/config"
	"github.com/tienhai2808/ecom_go/internal/handler"
	"github.com/tienhai2808/ecom_go/internal/repository"
	"github.com/tienhai2808/ecom_go/internal/security"
	"github.com/gin-gonic/gin"
)

func NewProductRouter(rg *gin.RouterGroup, cfg *config.Config, userRepo repository.UserRepository, productHdl *handler.ProductHandler) {
	accessName := cfg.App.AccessName
	secretKey := cfg.App.JWTSecret

	product := rg.Group("/products")
	{
		product.GET("", productHdl.GetAllProducts)

		product.GET("/:id", security.RequireAuth(accessName, secretKey, userRepo), security.RequireAdmin(), productHdl.GetProductByID)

		product.POST("", security.RequireAuth(accessName, secretKey, userRepo), security.RequireAdmin(), productHdl.CreateProduct)

		product.PATCH("/:id", security.RequireAuth(accessName, secretKey, userRepo), security.RequireAdmin(), productHdl.UpdateProduct)

		product.DELETE("/:id", security.RequireAuth(accessName, secretKey, userRepo), security.RequireAdmin(), productHdl.DeleteProduct)

		product.DELETE("", security.RequireAuth(accessName, secretKey, userRepo), security.RequireAdmin(), productHdl.DeleteProducts)
	}
}
