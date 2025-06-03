package model

import "time"

type Product struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Brand       string    `gorm:"type:varchar(50);not null" json:"brand"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Inventory   int       `gorm:"type:int;not null" json:"inventory"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`
}
