package request

type AddAddressRequest struct {
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required,len=10"`
	Address     string `json:"address" binding:"required"`
	Commune     string `json:"commune" binding:"required"`
	Province    string `json:"province" binding:"required"`
	IsDefault   *bool  `json:"is_default" binding:"required"`
}

type UpdateAddressRequest struct {
	FullName    *string `json:"full_name" binding:"omitempty"`
	PhoneNumber *string `json:"phone_number" binding:"omitempty,len=10"`
	Address     *string `json:"address" binding:"omitempty"`
	Commune     *string `json:"commune" binding:"omitempty,min=4"`
	Province    *string `json:"province" binding:"omitempty,min=4"`
	IsDefault   *bool   `json:"is_default" binding:"omitempty"`
}
