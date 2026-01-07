package dto

type ItemDTO struct {
	Name        string `json:"name" form:"name" validate:"required"`
	Description string `json:"description" form:"description" validate:"required"`
	Category    string `json:"category" form:"category" validate:"required"`
	PricePerDay int    `json:"price_per_day" form:"price_per_day" validate:"required"`
	Status      string `json:"status" form:"status" validate:"required"`
}
