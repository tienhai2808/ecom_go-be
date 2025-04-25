package main

import (
	"e-commerce/internal/auth"
	"e-commerce/internal/cache"
	"e-commerce/internal/common"
	"e-commerce/internal/config"
	"e-commerce/internal/database"
	"e-commerce/internal/product"
	"e-commerce/internal/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	if err := common.InitValidator(); err != nil {
		log.Fatal("Failed to init validator:", err)
	}
	r := gin.Default()

	appConfig, err := config.LoadAppConfig(); 
	if err != nil {
		log.Fatalf("❤️ Load cấu hình thất bại: %v", err)
	}

	// Kết nối MySQL
	db, err := database.ConnectToDatabase()
	if err != nil {
		log.Fatalf("❤️ Không thể kết nối database: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❤️ Không thể kết nối cơ sở dữ liệu thô: %v", err)
	}
	defer sqlDB.Close()
	fmt.Println("💚 Kết nối MySQL thành công")

	// Kết nối Redis
	redisClient, err := cache.ConnectToRedis()
	if err != nil {
		log.Fatalf("Không thể kết nối Redis: %v", err)
	}
	defer redisClient.Close()
	fmt.Println("💚 Kết nối Redis thành công")

	appCtx := &common.AppContext{
		DB:          db,
		Redis:       redisClient,
		Config:      appConfig,
	}

	api := r.Group(appConfig.App.ApiPrefix)
	{
		auth.AuthRouter(api, appCtx)
		user.UserRouter(api, appCtx)
		product.ProductRouter(api, appCtx)
	}

	r.Run(appConfig.App.Host + ":" + appConfig.App.Port)
}
