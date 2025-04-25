package database

import (
	"e-commerce/internal/config"
	"e-commerce/internal/user"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Hàm kết nối với MySQL sử dụng thông tin cấu hình
func ConnectToDatabase() (*gorm.DB, error) {
	// Lấy cấu hình database
	cfg, err := config.LoadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	// Tạo chuỗi kết nối MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	// Kết nối tới MySQL
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("❤️ failed to connect to database: %v", err)
	}

	err = RunAutoMigrations(db)
	if err != nil {
		return nil, fmt.Errorf("❤️ failed to run auto migrations: %v", err)
	}

	return db, nil
}

func RunAutoMigrations(db *gorm.DB) error {
	err := db.AutoMigrate(&user.User{}, &user.Profile{})
	if err != nil {
		return fmt.Errorf("❤️ failed to migrate: %v", err)
	}

	return nil
}
