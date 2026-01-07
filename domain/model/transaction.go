package model

import "time"

type TransactionStatus string
type BookingLifecycleStatus string

const (
	TransactionPendingPayment TransactionStatus = "PENDING_PAYMENT"
	TransactionPaid           TransactionStatus = "PAID"
	TransactionOngoing        TransactionStatus = "ONGOING"
	TransactionCompleted      TransactionStatus = "COMPLETED"

	BookingInProgress BookingLifecycleStatus = "IN_PROGRESS"
	BookingFinished   BookingLifecycleStatus = "FINISHED"
)

type Transaction struct {
	IDTransaction int `json:"id_transaction" gorm:"column:id_transaction;primaryKey;autoIncrement"`

	OrderID string `json:"order_id" gorm:"column:order_id;type:varchar(50);unique;not null"`

	ItemID   int `json:"item_id" gorm:"column:item_id;not null"`
	RenterID int `json:"renter_id" gorm:"column:renter_id;not null"`

	StartDate time.Time `json:"start_date" gorm:"column:start_date;not null"`
	EndDate   time.Time `json:"end_date" gorm:"column:end_date;not null"`

	TotalPrice             int64                  `json:"total_price" gorm:"column:total_price;not null"`
	TransactionStatus      TransactionStatus      `json:"transaction_status" gorm:"column:transaction_status;type:varchar(30);default:'PENDING_PAYMENT'"`
	BookingLifecycleStatus BookingLifecycleStatus `json:"booking_lifecycle_status" gorm:"column:booking_lifecycle_status;type:varchar(20);default:'IN_PROGRESS'"`

	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`

	Payment *Payment `json:"payment" gorm:"foreignKey:TransactionID;references:IDTransaction"`

	Item   Item `json:"item" gorm:"foreignKey:ItemID;references:Item_ID"`
	Renter User `json:"renter" gorm:"foreignKey:RenterID;references:Iduser"`
}

func (Transaction) TableName() string {
	return "transactions"
}
