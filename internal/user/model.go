package user

import (
	"time"
)

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

type Address struct {
	ID string `gorm:"type:char(36);primaryKey" json:"-"`
	FirstName string     `gorm:"type:varchar(255);nullable" json:"first_name"`
	LastName    string     `gorm:"type:varchar(255);nullable" json:"last_name"`
	PhoneNumber string     `gorm:"type:varchar(20);not null" json:"phone_number"`
	Address string `gorm:"type:varchar(255);not null" json:"address"`
	Commune string `gorm:"type:varchar(255);not null" json:"commune"`
	District string `gorm:"type:varchar(255);not null" json:"district"`
	Province string `gorm:"type:varchar(255);not null" json:"province"`
	IsDefault bool `gorm:"type:boolean;not null" json:"is_default"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"-"`

	UserID string `gorm:"type:char(36);not null" json:"-"`
}

type Profile struct {
	ID          string     `gorm:"type:char(36);primaryKey" json:"-"`
	FirstName   string     `gorm:"type:varchar(255);nullable" json:"first_name"`
	LastName    string     `gorm:"type:varchar(255);nullable" json:"last_name"`
	PhoneNumber string     `gorm:"type:varchar(20);nullable" json:"phone_number"`
	DOB         time.Time  `gorm:"type:date;default:'2000-01-01 00:00:00.000'" json:"dob"`
	Gender      UserGender `gorm:"type:enum('male','female','other');default:'other'" json:"gender"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"-"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"-"`

	UserID string `gorm:"type:char(36);not null;unique" json:"-"`
}

type User struct {
	ID        string    `gorm:"type:char(36);primaryKey" json:"id"`
	Username  string    `gorm:"type:varchar(50);not null;uniqueIndex:idx_username_email" json:"username"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_username_email" json:"email"`
	Role      UserRole  `gorm:"type:enum('user','admin', 'contributor', 'staff');default:'user';not null" json:"role"`
	Password  string    `gorm:"type:varchar(512);not null" json:"-"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"-"`

	Profile Profile `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"profile"`
	Address []Address `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
