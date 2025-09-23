package model

import "time"

type Product struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Brand       string    `gorm:"type:varchar(50);not null" json:"brand"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Inventory   uint      `gorm:"type:int;not null" json:"inventory"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`
	CategoryID  string    `gorm:"type:char(36)" json:"-"`

	Category *Category `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category"`
	Images   []*Image  `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"images"`
}
