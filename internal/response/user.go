package response

import "time"

type AuthResponse struct {
	ID        string          `json:"id"`
	Username  string          `json:"username"`
	Email     string          `json:"email"`
	CreatedAt time.Time       `json:"created_at"`
	Profile   ProfileResponse `json:"profile"`
}

type ProfileResponse struct {
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `json:"phone_number"`
	DOB         time.Time `json:"dob"`
	Gender      string    `json:"gender"`
}
