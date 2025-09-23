package request

import "time"

type UpdateProfileRequest struct {
	FirstName   *string    `json:"first_name,omitempty"`
	LastName    *string    `json:"last_name,omitempty"`
	Gender      *string     `json:"gender" binding:"omitempty,oneof=other male female"`
	DOB         *time.Time `json:"dob,omitempty"`
	PhoneNumber *string    `json:"phone_number" binding:"omitempty,len=10"`
}
