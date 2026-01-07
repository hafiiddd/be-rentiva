package repository

import (
	"back-end/domain/model"

	"gorm.io/gorm"
)

type dbTransaction struct {
	conn *gorm.DB
}

func NewTransactionRepository(conn *gorm.DB) TransactionRepository {
	return &dbTransaction{conn: conn}
}

type TransactionRepository interface {
	Create(tx model.Transaction) (model.Transaction, error)
	FindByOrderID(orderID string) (model.Transaction, error)
	FindByID(id int) (model.Transaction, error)
	UpdateTransactionStatus(orderID string, status model.TransactionStatus) error
	UpdateBookingLifecycle(orderID string, status model.BookingLifecycleStatus) error
	FindByUserAndStatuses(userID int, statuses []model.TransactionStatus) ([]model.Transaction, error)
}

func (r *dbTransaction) Create(tx model.Transaction) (model.Transaction, error) {
	if err := r.conn.Create(&tx).Error; err != nil {
		return model.Transaction{}, err
	}
	return tx, nil
}

func (r *dbTransaction) FindByOrderID(orderID string) (model.Transaction, error) {
	var tx model.Transaction
	if err := r.conn.
		Preload("Payment").
		Preload("Item").
		Preload("Item.Owner").
		Preload("Renter").
		Where("order_id = ?", orderID).
		First(&tx).Error; err != nil {
		return tx, err
	}
	return tx, nil
}

func (r *dbTransaction) FindByID(id int) (model.Transaction, error) {
	var tx model.Transaction
	if err := r.conn.
		Preload("Item").
		Preload("Payment").
		Where("id_transaction = ?", id).
		First(&tx).Error; err != nil {
		return tx, err
	}
	return tx, nil
}

func (r *dbTransaction) UpdateTransactionStatus(orderID string, status model.TransactionStatus) error {
	return r.conn.Model(&model.Transaction{}).
		Where("order_id = ?", orderID).
		Update("transaction_status", status).Error
}

func (r *dbTransaction) UpdateBookingLifecycle(orderID string, status model.BookingLifecycleStatus) error {
	return r.conn.Model(&model.Transaction{}).
		Where("order_id = ?", orderID).
		Update("booking_lifecycle_status", status).Error
}

// FindByUserAndStatuses mengambil transaksi berdasarkan status dan kepemilikan (renter/owner).
func (r *dbTransaction) FindByUserAndStatuses(userID int, statuses []model.TransactionStatus) ([]model.Transaction, error) {
	var txs []model.Transaction
	err := r.conn.
		Joins("JOIN items ON items.item_id = transactions.item_id").
		Joins("JOIN payments ON payments.order_id = transactions.order_id").
		Where("transactions.transaction_status IN ?", statuses).
		Where(r.conn.Where("transactions.renter_id = ?", userID).Or("items.owner_id = ?", userID)).
		Preload("Item").
		Preload("Item.Owner").
		Preload("Renter").
		Preload("Payment").
		Find(&txs).Error
	if err != nil {
		return nil, err
	}
	return txs, nil
}
