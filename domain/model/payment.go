package model

import "time"

type PaymentStatus string

const (
	PaymentPending PaymentStatus = "PENDING"
	PaymentPaid    PaymentStatus = "PAID"
	PaymentFailed  PaymentStatus = "FAILED"
	PaymentExpired PaymentStatus = "EXPIRED"
)

type Payment struct {
	IDPayment int `json:"id_payment" gorm:"column:id_payment;primaryKey;autoIncrement"`

	TransactionID int    `json:"transaction_id" gorm:"column:transaction_id;not null"`
	OrderID       string `json:"order_id" gorm:"column:order_id;type:varchar(50);unique;not null"`

	PaymentType   string        `json:"payment_type" gorm:"column:payment_type;type:varchar(50)"`
	PaymentMethod string        `json:"payment_method" gorm:"column:payment_method;type:varchar(50)"`
	PaymentStatus PaymentStatus `json:"payment_status" gorm:"column:payment_status;type:varchar(20);default:'PENDING'"`
	PaymentRef    string        `json:"payment_ref" gorm:"column:payment_ref;type:text"`
	VANumber      string        `json:"va_number" gorm:"column:va_number;type:varchar(50)"`
	ExpiryTime    *time.Time    `json:"expiry_time" gorm:"column:expiry_time"`

	PaidAt *time.Time `json:"paid_at" gorm:"column:paid_at"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (Payment) TableName() string {
	return "payments"
}
