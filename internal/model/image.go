package model

import "time"

type Image struct {
	ID        int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Url       string    `gorm:"type:varchar(255)" json:"url"`
	PublicID  string    `gorm:"type:varchar(50)" json:"public_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	ProductID int64     `gorm:"type:bigint;not null" json:"product_id"`

	Product *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product"`
}
