package request

import (
	"backend/internal/model"
	"time"
)

type CreateUserRequest struct {
	Username    string           `json:"username" binding:"required,min=3"`
	Email       string           `json:"email" binding:"required,email"`
	Password    string           `json:"password" binding:"required,min=6"`
	Role        model.UserRole   `json:"role" binding:"required,oneof=user admin contributor staff"`
	Gender      model.UserGender `json:"gender" binding:"required,oneof=other male female"`
	FirstName   string           `json:"first_name" binding:"omitempty"`
	LastName    string           `json:"last_name" binding:"omitempty"`
	PhoneNumber string           `json:"phone_number" binding:"omitempty,len=10"`
	DOB         time.Time        `json:"dob" binding:"omitempty"`
}

type UpdateUserRequest struct {
	Username    *string           `json:"username" binding:"omitempty,min=3"`
	Email       *string           `json:"email" binding:"omitempty,email"`
	Password    *string           `json:"password" binding:"omitempty,min=6"`
	Role        *model.UserRole   `json:"role" binding:"omitempty,oneof=user admin contributor staff"`
	Gender      *model.UserGender `json:"gender" binding:"omitempty,oneof=other male female"`
	FirstName   *string           `json:"first_name" binding:"omitempty"`
	LastName    *string           `json:"last_name" binding:"omitempty"`
	PhoneNumber *string           `json:"phone_number" binding:"omitempty,len=10"`
	DOB         *time.Time        `json:"dob" binding:"omitempty"`
}

type DeleteManyUsersRequest struct {
	UserIds []string `json:"user_ids" binding:"required,min=1,dive,uuid"`
}
