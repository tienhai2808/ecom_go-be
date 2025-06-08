package model

import "time"

type Image struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	FileName    string    `gorm:"type:varchar(255);not null" json:"file_name"`
	FileType    string    `gorm:"type:varchar(100);not null" json:"file_type"`
	Image       []byte    `gorm:"type:blob;not null" json:"image"`
	DownloadUrl string    `gorm:"type:varchar(512);not null" json:"download_url"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`

	ProductID string  `gorm:"type:char(36)" json:"-"`
	Product   Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
