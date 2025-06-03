package dto

type RegistrationData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Otp      string `json:"otp"`
	Attempts int    `json:"attempts"`
}