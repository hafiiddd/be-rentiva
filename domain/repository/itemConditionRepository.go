package repository

import (
	"back-end/domain/model"

	"gorm.io/gorm"
)

type ItemConditionRepository interface {
	Create(cond model.ItemCondition) (model.ItemCondition, error)
	ExistsByTransactionAndType(transactionID int, conditionType model.ConditionType) (bool, error)
}

type dbItemCondition struct {
	conn *gorm.DB
}

func NewItemConditionRepository(conn *gorm.DB) ItemConditionRepository {
	return &dbItemCondition{conn: conn}
}

func (r *dbItemCondition) Create(cond model.ItemCondition) (model.ItemCondition, error) {
	if err := r.conn.Create(&cond).Error; err != nil {
		return model.ItemCondition{}, err
	}
	return cond, nil
}

func (r *dbItemCondition) ExistsByTransactionAndType(transactionID int, conditionType model.ConditionType) (bool, error) {
	var count int64
	if err := r.conn.Model(&model.ItemCondition{}).
		Where("transaction_id = ? AND condition_type = ?", transactionID, conditionType).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
