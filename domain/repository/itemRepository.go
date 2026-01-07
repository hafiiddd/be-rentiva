package repository

import (
	"back-end/domain/model"

	"gorm.io/gorm"
)

type dbItem struct {
	conn *gorm.DB
}

// FindItemByid implements ItemRepository.
func (d *dbItem) FindItemByid(idItem int) (model.Item, error) {
	var item model.Item
	err := d.conn.Preload("Owner").Where("item_id = ?", idItem).Find(&item).Error
	if err != nil {
		return item, err
	}

	return item, nil
}

// FindItemByUserId implements ItemRepository.
func (d *dbItem) FindItemByUserId(id int) ([]model.Item, error) {
	var item []model.Item
	err := d.conn.Where("owner_id = ?", id).Find(&item).Error
	if err != nil {
		return item, err
	}

	return item, nil
}

// CreateItem implements ItemRepository.
func (d *dbItem) CreateItem(newItem model.Item) (model.Item, error) {
	if err := d.conn.Create(&newItem).Error; err != nil {
		return model.Item{}, err
	}

	// Muat relasi owner setelah insert
	if err := d.conn.Preload("Owner").First(&newItem, newItem.Item_ID).Error; err != nil {
		return model.Item{}, err
	}

	return newItem, nil
}

// UpdateItem implements ItemRepository.
func (d *dbItem) UpdateItem(id int, updated model.Item) (model.Item, error) {
	var existing model.Item
	if err := d.conn.First(&existing, id).Error; err != nil {
		return model.Item{}, err
	}

	updated.Item_ID = existing.Item_ID
	if err := d.conn.Model(&existing).Updates(updated).Error; err != nil {
		return model.Item{}, err
	}

	if err := d.conn.Preload("Owner").First(&existing, id).Error; err != nil {
		return model.Item{}, err
	}

	return existing, nil
}

// DeleteItem implements ItemRepository.
func (d *dbItem) DeleteItem(id int) error {
	return d.conn.Delete(&model.Item{}, id).Error
}

type ItemRepository interface {
	FindItemByUserId(id int) ([]model.Item, error)
	FindItemByid(idItemn int) (model.Item, error)
	CreateItem(newItem model.Item) (model.Item, error)
	UpdateItem(id int, updated model.Item) (model.Item, error)
	DeleteItem(id int) error
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &dbItem{conn: db}
}
