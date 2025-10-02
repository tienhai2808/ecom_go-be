package initialization

import (
	"database/sql"
	"fmt"
	"github.com/tienhai2808/ecom_go/config"
	"github.com/tienhai2808/ecom_go/internal/model"

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
	&model.Inventory{},
}

type DB struct {
	Gorm *gorm.DB
	sql  *sql.DB
}

func InitMySQL(cfg *config.Config) (*DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	gDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("kết nối MySQL thất bại: %w", err)
	}

	if err = runAutoMigrations(gDB); err != nil {
		return nil, fmt.Errorf("chuyển dịch DB thất bại: %w", err)
	}

	sqlDB, err := gDB.DB()
	if err != nil {
		return nil, fmt.Errorf("không lấy được sql.DB: %w", err)
	}

	return &DB{
		gDB,
		sqlDB,
	}, nil
}

func (d *DB) Close() {
	_ = d.sql.Close()
}

func runAutoMigrations(db *gorm.DB) error {
	return db.AutoMigrate(allModels...)
}
