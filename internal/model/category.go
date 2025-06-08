package model

import "time"

type Category struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(150);not null" json:"name"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`

	Products []Product `gorm:"foreignKey:CategoryID;references:ID" json:"products"`
}
