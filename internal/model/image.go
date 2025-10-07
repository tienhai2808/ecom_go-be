package model

type Image struct {
	ID          int64  `gorm:"type:bigint;primaryKey" json:"id"`
	Url         string `gorm:"type:varchar(255)" json:"url"`
	PublicID    string `gorm:"type:varchar(50)" json:"public_id"`
	IsThumbnail bool   `gorm:"type:boolean;not null" json:"is_thumbnail"`
	SortOrder   int    `gorm:"type:int;not null" json:"sort_order"`
	ProductID   int64  `gorm:"type:bigint;not null" json:"product_id"`

	Product *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"product"`
}
