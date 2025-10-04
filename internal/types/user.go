package types

import "time"

type UserData struct {
	ID        int64       `json:"id"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	Role      string      `json:"role"`
	CreatedAt time.Time   `json:"created_at"`
	Profile   ProfileData `json:"profile"`
}

type ProfileData struct {
	ID          int64      `json:"id"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	PhoneNumber string     `json:"phone_number"`
	DOB         *time.Time `json:"dob"`
	Gender      string     `json:"gender"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
