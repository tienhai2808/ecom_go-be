package model

import "time"

type Profile struct {
	ID          string    `gorm:"type:char(36);primaryKey" json:"-"`
	FirstName   string    `gorm:"type:varchar(255);nullable" json:"first_name"`
	LastName    string    `gorm:"type:varchar(255);nullable" json:"last_name"`
	PhoneNumber string    `gorm:"type:varchar(20);nullable" json:"phone_number"`
	DOB         time.Time `gorm:"type:date;default:'2000-01-01 00:00:00.000'" json:"dob"`
	Gender      string    `gorm:"type:enum('male','female','other');default:'other'" json:"gender"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"-"`
	UserID      string    `gorm:"type:char(36);not null;unique" json:"-"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
}
