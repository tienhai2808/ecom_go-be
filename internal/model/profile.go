package model

import "time"

type Profile struct {
	ID          int64      `gorm:"type:bigint;primaryKey" json:"id"`
	FirstName   string     `gorm:"type:varchar(255);nullable" json:"first_name"`
	LastName    string     `gorm:"type:varchar(255);nullable" json:"last_name"`
	PhoneNumber string     `gorm:"type:varchar(20);nullable" json:"phone_number"`
	DOB         *time.Time `gorm:"type:date" json:"dob"`
	Gender      string     `gorm:"type:enum('male','female','other');default:'other'" json:"gender"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	UserID      int64      `gorm:"type:bigint;not null;unique" json:"user_id"`

	User *User `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
}
