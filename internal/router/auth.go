package router

import (
	"backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(rg *gin.RouterGroup, authHandler handler.AuthHandler) {
	auth := rg.Group(("/auth"))
	{
		auth.POST("/signup", authHandler.Signup)
	}
}