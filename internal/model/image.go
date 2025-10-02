package model

import "time"

type Image struct {
	ID         int64     `gorm:"type:bigint;primaryKey" json:"id"`
	ImageUrl   string    `gorm:"type:varchar(255);not null" json:"image_url"`
	ImageKitID string    `gorm:"type:varchar(50);not null" json:"-"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"-"`
	ProductID  int64     `gorm:"type:bigint;not null" json:"-"`

	Product *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product"`
}
