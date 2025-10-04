package response

import "time"

type AddressResponse struct {
	ID          int64     `json:"id"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	Address     string    `json:"address"`
	Commune     string    `json:"commune"`
	Province    string    `json:"province"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BaseAddressResponse struct {
	ID          int64  `json:"id"`
	FullName    string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Commune     string `json:"commune"`
	Province    string `json:"province"`
}
