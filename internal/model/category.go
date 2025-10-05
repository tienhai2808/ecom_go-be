package model

import "time"

type Category struct {
	ID        int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(150);not null" json:"name"`
	Slug      string    `gorm:"type:varchar(150);not null;unique" json:"slug"`
	IsDeleted bool      `gorm:"type:boolean;not null;default:false" json:"is_deleted"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Products []*Product `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL" json:"products"`
}
