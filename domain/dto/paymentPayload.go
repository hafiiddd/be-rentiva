package dto

type SnapResultRequest struct {
	OrderID       string `json:"order_id" validate:"required"`
	PaymentType   string `json:"payment_type" validate:"required"`
	PaymentMethod string `json:"payment_method" validate:"required"`
	VANumber      string `json:"va_number" validate:"required"`
}
