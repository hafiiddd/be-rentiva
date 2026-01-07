package dto

type CreateTransactionRequest struct {
	ItemID     int    `json:"item_id" validate:"required"`
	StartDate  string `json:"start_date" validate:"required,datetime=2006-01-02"`
	EndDate    string `json:"end_date" validate:"required,datetime=2006-01-02"`
	TotalPrice int64  `json:"total_price" validate:"required"`
}

type MidtransWebhookRequest struct {
	TransactionStatus string `json:"transaction_status"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
	OrderID           string `json:"order_id"`
	TransactionID     string `json:"transaction_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
}
