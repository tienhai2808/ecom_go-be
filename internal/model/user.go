package model

import "time"

type User struct {
	ID        int64     `gorm:"type:bigint;primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null;unique" json:"username"`
	Email     string    `gorm:"type:varchar(255);not null;unique" json:"email"`
	Role      string    `gorm:"type:enum('user','admin');default:'user';not null" json:"role"`
	Password  string    `gorm:"type:varchar(512);not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`

	Profile *Profile   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile"`
	Address []*Address `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
