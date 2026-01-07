package service

import (
	"back-end/domain/model"
	"back-end/domain/repository"
	"back-end/helper"
	"net/http"
)

// itemService adalah implementasi konkret
type itemService struct {
	itemRepo repository.ItemRepository // Dependensi ke repository
}

// GetItemByid implements ItemService.
func (s *itemService) GetItemByid(idItemn int) helper.Response {
	items, err := s.itemRepo.FindItemByid(idItemn)
	if err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal mengambil data barang di database",
			Error:    err.Error(),
		}
	}

	return helper.Response{
		Status:   http.StatusCreated,
		Messages: "Berhasil mengambil data barang",
		Data:     items,
	}
}

type ItemService interface {
	GetItemsByUserID(userId int) helper.Response
	GetItemByid(idItemn int) helper.Response
	CreateItem(newItem model.Item) helper.Response
	UpdateItem(id int, updated model.Item) helper.Response
	DeleteItem(id int) helper.Response
}

// NewItemService membuat instance service baru
func NewItemService(itemRepo repository.ItemRepository) ItemService {
	return &itemService{itemRepo: itemRepo}
}

// GetItemsByUserID (Implementasi)
func (s *itemService) GetItemsByUserID(userId int) helper.Response {
	// 1. Panggil repository
	items, err := s.itemRepo.FindItemByUserId(userId)
	if err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal mengambil data barang di database",
			Error:    err.Error(),
		}
	}

	return helper.Response{
		Status:   http.StatusCreated,
		Messages: "Berhasil mengambil data barang",
		Data:     items,
	}
}

// CreateItem (Implementasi)
func (s *itemService) CreateItem(newItem model.Item) helper.Response {
	item, err := s.itemRepo.CreateItem(newItem)
	if err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal menyimpan data barang",
			Error:    err.Error(),
		}
	}

	return helper.Response{
		Status:   http.StatusCreated,
		Messages: "Berhasil menambahkan barang",
		Data:     item,
	}
}

// UpdateItem (Implementasi)
func (s *itemService) UpdateItem(id int, updated model.Item) helper.Response {
	item, err := s.itemRepo.UpdateItem(id, updated)
	if err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal mengupdate data barang",
			Error:    err.Error(),
		}
	}

	return helper.Response{
		Status:   http.StatusOK,
		Messages: "Berhasil mengupdate barang",
		Data:     item,
	}
}

// DeleteItem (Implementasi)
func (s *itemService) DeleteItem(id int) helper.Response {
	if err := s.itemRepo.DeleteItem(id); err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal menghapus data barang",
			Error:    err.Error(),
		}
	}

	return helper.Response{
		Status:   http.StatusOK,
		Messages: "Berhasil menghapus barang",
	}
}
