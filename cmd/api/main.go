package main

import (
	"e-commerce/internal/auth"
	"e-commerce/internal/cache"
	"e-commerce/internal/common"
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

	// Kết nối đến MySQL
	db, err := database.ConnectToDatabase()
	if err != nil {
		log.Fatalf("❤️ Could not connect to database: %v", err)
	}

	// Kiểm tra kết nối
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("❤️ Could not get raw database connection: %v", err)
	}
	defer sqlDB.Close()

	fmt.Println("💚 Successfully connected to the database")

	redisClient, err := cache.ConnectToRedis()
	if err != nil {
		log.Fatalf("Không thể kết nối Redis: %v", err)
	}
	defer redisClient.Close()
	fmt.Println("💚 Kết nối Redis thành công")

	fmt.Println("💚 Successfully connected to the cache")

	api := r.Group("/e-commerce")
	{
		auth.AuthRouter(api)
		user.UserRouter(api)
		product.ProductRouter(api)
	}

	r.Run(":8080")
}
