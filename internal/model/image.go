package model

import "time"

type Image struct {
	ID         string    `gorm:"type:char(36);primaryKey" json:"id"`
	ImageUrl   string    `gorm:"type:varchar(255);not null" json:"image_url"`
	ImageKitID string    `gorm:"type:varchar(50);not null" json:"-"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"-"`
	ProductID  string    `gorm:"type:char(36);not null" json:"-"`

	Product *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product"`
}
