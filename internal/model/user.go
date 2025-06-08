package model

import "time"

type UserRole string

type UserGender string

const (
	RoleUser        UserRole = "user"
	RoleAdmin       UserRole = "admin"
	RoleContributor UserRole = "contributor"
	RoleStaff       UserRole = "staff"
)

const (
	GenderMale   UserGender = "male"
	GenderFemale UserGender = "female"
	GenderOther  UserGender = "other"
)

type User struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_username_email" json:"username"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_username_email" json:"email"`
	Role      UserRole  `gorm:"type:enum('user','admin', 'contributor', 'staff');default:'user';not null" json:"role"`
	Password  string    `gorm:"type:varchar(512);not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`

	Profile Profile   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile"`
	Address []Address `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
