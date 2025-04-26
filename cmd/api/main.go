package main

import (
	"backend/internal/auth"
	"backend/internal/cache"
	"backend/internal/common"
	"backend/internal/config"
	"backend/internal/database"
	"backend/internal/product"
	"backend/internal/user"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
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
