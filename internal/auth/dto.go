package auth

// SignupRequest defines the structure for signup request
type SignupRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// VerifySignupRequest defines the structure for signup verification
type VerifySignupRequest struct {
	RegistrationToken string `json:"registration_token" binding:"required,uuid4"`
	Otp              string `json:"otp" binding:"required,len=6,numeric"`
}

// SigninRequest defines the structure for signin request
type SigninRequest struct {
	Username string `json:"username" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegistrationData holds data during registration process
type RegistrationData struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	OTP      string `json:"otp"`
	Attempts int    `json:"attempts"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Profile  interface{} `json:"profile,omitempty"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}