package model

import "time"

type Product struct {
	ID          int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string    `gorm:"type:varchar(255);not null;unique" json:"brand"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`
	CategoryID  int64     `gorm:"type:bigint" json:"-"`

	Category  *Category  `gorm:"foreignKey:CategoryID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"category"`
	Images    []*Image   `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"images"`
	Inventory *Inventory `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE;OnDelete:CASCADE" json:"inventory"`
}
