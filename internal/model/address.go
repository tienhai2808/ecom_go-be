package model

import "time"

type Address struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"id"`
	FirstName   string    `gorm:"type:varchar(255);nullable" json:"first_name"`
	LastName    string    `gorm:"type:varchar(255);nullable" json:"last_name"`
	PhoneNumber string    `gorm:"type:varchar(20);not null" json:"phone_number"`
	Address     string    `gorm:"type:varchar(255);not null" json:"address"`
	Commune     string    `gorm:"type:varchar(255);not null" json:"commune"`
	District    string    `gorm:"type:varchar(255);not null" json:"district"`
	Province    string    `gorm:"type:varchar(255);not null" json:"province"`
	IsDefault   bool      `gorm:"type:boolean;not null" json:"is_default"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`

	UserID string `gorm:"type:char(36);not null" json:"-"`
}