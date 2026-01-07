package dto

import (
	"back-end/domain/model"
	"time"
)

type BookingItemInfo struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Photo string `json:"photo"`
}

type BookingItemResponse struct {
	OrderID                string                       `json:"order_id"`
	TransactionStatus      model.TransactionStatus      `json:"transaction_status"`
	BookingLifecycleStatus model.BookingLifecycleStatus `json:"booking_lifecycle_status"`
	StartDate              time.Time                    `json:"start_date"`
	EndDate                time.Time                    `json:"end_date"`
	TotalPrice             int64                        `json:"total_price"`
	Item                   BookingItemInfo              `json:"item"`
}
