package dto

type AuthRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ResisterDTO struct {

	// Data dari Langkah 1
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`

	// Data dari Langkah 4
	FullName string `json:"full_name" validate:"required"`
	Nik      string `json:"nik" validate:"required"`
	Ttl      string `json:"ttl" validate:"required"`
	Address  string `json:"address" validate:"required"`
}
