package model

import "time"

type Address struct {
	ID          int64     `gorm:"type:bigint;primaryKey" json:"id"`
	FullName    string    `gorm:"type:varchar(255)" json:"full_name"`
	PhoneNumber string    `gorm:"type:varchar(20);not null" json:"phone_number"`
	Address     string    `gorm:"type:varchar(255);not null" json:"address"`
	Commune     string    `gorm:"type:varchar(255);not null" json:"commune"`
	Province    string    `gorm:"type:varchar(255);not null" json:"province"`
	IsDefault   bool      `gorm:"type:boolean;not null" json:"is_default"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`
	UserID      int64     `gorm:"type:bigint;not null" json:"-"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
