package service

import (
	"back-end/domain/dto"
	"back-end/domain/model"
	"back-end/domain/repository"
	"back-end/helper"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type UploadConditionFile struct {
	Filename    string
	Size        int64
	ContentType string
	Reader      io.Reader
}

type ItemConditionService interface {
	UploadItemCondition(ctx context.Context, orderID string, userID int, req dto.UploadItemConditionRequest, file UploadConditionFile) helper.Response
}

type itemConditionService struct {
	txRepo          repository.TransactionRepository
	itemCondRepo    repository.ItemConditionRepository
	storageService  StorageService
}

func NewItemConditionService(
	txRepo repository.TransactionRepository,
	itemCondRepo repository.ItemConditionRepository,
	storageService StorageService,
) ItemConditionService {
	return &itemConditionService{
		txRepo:         txRepo,
		itemCondRepo:   itemCondRepo,
		storageService: storageService,
	}
}

func (s *itemConditionService) UploadItemCondition(ctx context.Context, orderID string, userID int, req dto.UploadItemConditionRequest, file UploadConditionFile) helper.Response {
	if !req.ConditionType.IsValid() {
		return helper.Response{Status: http.StatusBadRequest, Messages: "condition_type tidak valid"}
	}
	if file.Reader == nil || file.Size == 0 {
		return helper.Response{Status: http.StatusBadRequest, Messages: "photo wajib diisi"}
	}

	tx, err := s.txRepo.FindByOrderID(orderID)
	if err != nil {
		return helper.Response{Status: http.StatusBadRequest, Messages: "Transaksi tidak ditemukan", Error: err.Error()}
	}
	if tx.BookingLifecycleStatus != model.BookingInProgress {
		return helper.Response{Status: http.StatusBadRequest, Messages: "Booking harus IN_PROGRESS"}
	}

	ownerID := int(tx.Item.OwnerID)
	isOwner := userID == ownerID
	isRenter := userID == tx.RenterID
	if !isOwner && !isRenter {
		return helper.Response{Status: http.StatusForbidden, Messages: "Anda bukan bagian dari transaksi ini"}
	}

	if resp := s.validateConditionRules(req.ConditionType, tx.TransactionStatus, isOwner, isRenter); resp != nil {
		return *resp
	}

	exists, err := s.itemCondRepo.ExistsByTransactionAndType(tx.IDTransaction, req.ConditionType)
	if err != nil {
		return helper.Response{Status: http.StatusInternalServerError, Messages: "Gagal cek duplikasi kondisi", Error: err.Error()}
	}
	if exists {
		return helper.Response{Status: http.StatusBadRequest, Messages: "Kondisi dengan tipe ini sudah diupload untuk transaksi ini"}
	}

	objectName := UniqueObjectName(
		fmt.Sprintf("item-conditions/%s/%s", orderID, strings.ToLower(string(req.ConditionType))),
		file.Filename,
	)
	key, err := s.storageService.Upload(ctx, objectName, file.Reader, file.Size, file.ContentType)
	if err != nil {
		return helper.Response{Status: http.StatusInternalServerError, Messages: "Gagal upload ke storage", Error: err.Error()}
	}

	photoURL := key
	if url, err := s.storageService.BuildPublicURL(key); err == nil {
		photoURL = url
	} else if url, err := s.storageService.GetPresignedURL(key); err == nil {
		photoURL = url
	}

	condition := model.ItemCondition{
		TransactionID: tx.IDTransaction,
		UserID:        userID,
		ConditionType: req.ConditionType,
		PhotoURL:      photoURL,
		Note:          req.Note,
	}
	condition, err = s.itemCondRepo.Create(condition)
	if err != nil {
		return helper.Response{Status: http.StatusInternalServerError, Messages: "Gagal menyimpan kondisi barang", Error: err.Error()}
	}

	if err := s.updateTransactionStatus(orderID, req.ConditionType); err != nil {
		return helper.Response{Status: http.StatusInternalServerError, Messages: "Gagal memperbarui status transaksi", Error: err.Error()}
	}

	return helper.Response{
		Status:   http.StatusCreated,
		Messages: "Berhasil upload kondisi barang",
		Data:     condition,
	}
}

func (s *itemConditionService) validateConditionRules(condType model.ConditionType, txStatus model.TransactionStatus, isOwner, isRenter bool) *helper.Response {
	switch condType {
	case model.ConditionBeforeSend:
		if !isOwner {
			return &helper.Response{Status: http.StatusForbidden, Messages: "BEFORE_SEND hanya boleh oleh OWNER"}
		}
		if txStatus != model.TransactionPaid {
			return &helper.Response{Status: http.StatusBadRequest, Messages: "Transaksi harus PAID untuk BEFORE_SEND"}
		}
	case model.ConditionAfterReceive:
		if !isRenter {
			return &helper.Response{Status: http.StatusForbidden, Messages: "AFTER_RECEIVE hanya boleh oleh RENTER"}
		}
		if txStatus != model.TransactionOngoing {
			return &helper.Response{Status: http.StatusBadRequest, Messages: "Transaksi harus ONGOING untuk AFTER_RECEIVE"}
		}
	case model.ConditionBeforeReturn:
		if !isRenter {
			return &helper.Response{Status: http.StatusForbidden, Messages: "BEFORE_RETURN hanya boleh oleh RENTER"}
		}
		if txStatus != model.TransactionOngoing {
			return &helper.Response{Status: http.StatusBadRequest, Messages: "Transaksi harus ONGOING untuk BEFORE_RETURN"}
		}
	case model.ConditionAfterReturn:
		if !isOwner {
			return &helper.Response{Status: http.StatusForbidden, Messages: "AFTER_RETURN hanya boleh oleh OWNER"}
		}
		if txStatus != model.TransactionOngoing {
			return &helper.Response{Status: http.StatusBadRequest, Messages: "Transaksi harus ONGOING untuk AFTER_RETURN"}
		}
	default:
		return &helper.Response{Status: http.StatusBadRequest, Messages: "condition_type tidak dikenal"}
	}
	return nil
}

func (s *itemConditionService) updateTransactionStatus(orderID string, condType model.ConditionType) error {
	switch condType {
	case model.ConditionBeforeSend:
		return s.txRepo.UpdateTransactionStatus(orderID, model.TransactionOngoing)
	case model.ConditionAfterReturn:
		return s.txRepo.UpdateTransactionStatus(orderID, model.TransactionCompleted)
	default:
		return nil
	}
}
