package database

import (
	"backend/internal/config"
	"backend/internal/model"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var allModels = []interface{}{
	&model.User{},
	&model.Profile{},
	&model.Address{},
	&model.Product{},
	&model.Category{},
	&model.Product{},
	&model.Image{},
}

func ConnectToDatabase(cfg *config.AppConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

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
	err := db.AutoMigrate(allModels...)
	if err != nil {
		return fmt.Errorf("❤️ failed to migrate: %v", err)
	}

	return nil
}
