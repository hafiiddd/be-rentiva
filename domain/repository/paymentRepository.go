package repository

import (
	"back-end/domain/model"

	"gorm.io/gorm"
)

type dbPayment struct {
	conn *gorm.DB
}

func NewPaymentRepository(conn *gorm.DB) PaymentRepository {
	return &dbPayment{conn: conn}
}

type PaymentRepository interface {
	Create(p model.Payment) (model.Payment, error)
	UpdateByOrderID(orderID string, updates map[string]interface{}) (model.Payment, error)
	UpdateSnapInfoByOrderID(orderID string, updates map[string]interface{}) (model.Payment, error)
	UpdateStatus(orderID string, status model.PaymentStatus, updates map[string]interface{}) (model.Payment, error)
}

func (r *dbPayment) Create(p model.Payment) (model.Payment, error) {
	if err := r.conn.Create(&p).Error; err != nil {
		return model.Payment{}, err
	}
	return p, nil
}

func (r *dbPayment) UpdateByOrderID(orderID string, updates map[string]interface{}) (model.Payment, error) {
	var payment model.Payment
	if err := r.conn.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		return payment, err
	}

	if err := r.conn.Model(&payment).Updates(updates).Error; err != nil {
		return payment, err
	}

	if err := r.conn.First(&payment, payment.IDPayment).Error; err != nil {
		return payment, err
	}

	return payment, nil
}

func (r *dbPayment) UpdateStatus(orderID string, status model.PaymentStatus, updates map[string]interface{}) (model.Payment, error) {
	updates["payment_status"] = status
	return r.UpdateByOrderID(orderID, updates)
}

func (r *dbPayment) UpdateSnapInfoByOrderID(orderID string, updates map[string]interface{}) (model.Payment, error) {
	var payment model.Payment
	if err := r.conn.Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		return payment, err
	}

	// Jangan ubah payment_status di sini
	delete(updates, "payment_status")

	if err := r.conn.Model(&payment).Updates(updates).Error; err != nil {
		return payment, err
	}

	if err := r.conn.First(&payment, payment.IDPayment).Error; err != nil {
		return payment, err
	}

	return payment, nil
}
