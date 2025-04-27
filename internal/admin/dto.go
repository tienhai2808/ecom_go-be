package admin

import (
	"backend/internal/user"
	"time"
)

type CreateUserRequest struct {
	Username    string          `json:"username" binding:"required,min=3"`
	Email       string          `json:"email" binding:"required,email"`
	Password    string          `json:"password" binding:"required,min=6"`
	Role        user.UserRole   `json:"role" binding:"required,oneof=user admin contributor staff"`
	Gender      user.UserGender `json:"gender" binding:"required,oneof=other male female"`
	FirstName   string          `json:"first_name"`
	LastName    string          `json:"last_name"`
	PhoneNumber string          `json:"phone_number"`
	DOB         time.Time       `json:"dob"`
}
