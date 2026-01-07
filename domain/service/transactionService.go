package service

import (
	"back-end/config"
	"back-end/domain/dto"
	"back-end/domain/model"
	"back-end/domain/repository"
	"back-end/helper"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type TransactionService interface {
	CreateTransaction(req dto.CreateTransactionRequest, renterID int) helper.Response
	HandleWebhook(req dto.MidtransWebhookRequest) helper.Response
	ListPaidTransactions(userID int) helper.Response
	SetOngoing(orderID string) helper.Response
	SetCompleted(orderID string) helper.Response
	FinishBooking(orderID string) helper.Response
}

type transactionService struct {
	txRepo         repository.TransactionRepository
	paymentRepo    repository.PaymentRepository
	paymentGateway PaymentGateway
	midtransConfig config.MidtransConfig
}

func NewTransactionService(
	txRepo repository.TransactionRepository,
	paymentRepo repository.PaymentRepository,
	paymentGateway PaymentGateway,
	midConfig config.MidtransConfig,
) TransactionService {
	return &transactionService{
		txRepo:         txRepo,
		paymentRepo:    paymentRepo,
		paymentGateway: paymentGateway,
		midtransConfig: midConfig,
	}
}

func (s *transactionService) CreateTransaction(req dto.CreateTransactionRequest, renterID int) helper.Response {
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return helper.Response{
			Status:   http.StatusBadRequest,
			Messages: "start_date tidak valid, gunakan format 2006-01-02",
			Error:    err.Error(),
		}
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return helper.Response{
			Status:   http.StatusBadRequest,
			Messages: "end_date tidak valid, gunakan format 2006-01-02",
			Error:    err.Error(),
		}
	}

	orderID := fmt.Sprintf("RENTIVA-%d", time.Now().UnixNano())

	tx := model.Transaction{
		OrderID:                orderID,
		ItemID:                 req.ItemID,
		RenterID:               renterID,
		StartDate:              startDate,
		EndDate:                endDate,
		TotalPrice:             req.TotalPrice,
		TransactionStatus:      model.TransactionPendingPayment,
		BookingLifecycleStatus: model.BookingInProgress,
	}

	tx, err = s.txRepo.Create(tx)
	if err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal membuat transaksi",
			Error:    err.Error(),
		}
	}

	_, err = s.paymentRepo.Create(model.Payment{
		TransactionID: tx.IDTransaction,
		OrderID:       orderID,
		PaymentStatus: model.PaymentPending,
	})
	if err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal membuat payment",
			Error:    err.Error(),
		}
	}

	snapToken, err := s.paymentGateway.CreateSnap(orderID, req.TotalPrice)
	if err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal membuat Snap token",
			Error:    err.Error(),
		}
	}

	return helper.Response{
		Status:   http.StatusCreated,
		Messages: "Transaksi berhasil dibuat",
		Data: map[string]interface{}{
			"order_id":    orderID,
			"snap_token":  snapToken,
			"booking":     tx,
			"payment_url": fmt.Sprintf("%s/snap/v2/vtweb/%s", s.midtransConfig.BaseUri, snapToken),
		},
	}
}

func (s *transactionService) HandleWebhook(req dto.MidtransWebhookRequest) helper.Response {
	if !s.verifySignature(req) {
		return helper.Response{
			Status:   http.StatusUnauthorized,
			Messages: "signature tidak valid",
		}
	}

	paymentStatus := mapMidtransStatus(req.TransactionStatus, req.FraudStatus)
	updates := map[string]interface{}{
		"payment_type":   req.PaymentType,
		"payment_method": req.PaymentType,
		"payment_ref":    req.TransactionID,
	}
	if paymentStatus == model.PaymentPaid {
		now := time.Now()
		updates["paid_at"] = &now
	}

	payment, err := s.paymentRepo.UpdateStatus(req.OrderID, paymentStatus, updates)
	if err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal memperbarui payment",
			Error:    err.Error(),
		}
	}

	switch paymentStatus {
	case model.PaymentPaid:
		_ = s.txRepo.UpdateTransactionStatus(req.OrderID, model.TransactionPaid)
		_ = s.txRepo.UpdateBookingLifecycle(req.OrderID, model.BookingInProgress)
	case model.PaymentFailed, model.PaymentExpired:
		_ = s.txRepo.UpdateTransactionStatus(req.OrderID, model.TransactionPendingPayment)
	}

	return helper.Response{
		Status:   http.StatusOK,
		Messages: "Webhook diproses",
		Data:     payment,
	}
}

func (s *transactionService) ListPaidTransactions(userID int) helper.Response {
	statuses := []model.TransactionStatus{
		model.TransactionPaid,
		model.TransactionOngoing,
		model.TransactionCompleted,
	}
	txs, err := s.txRepo.FindByUserAndStatuses(userID, statuses)
	if err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal mengambil transaksi",
			Error:    err.Error(),
		}
	}

	resp := make([]dto.BookingItemResponse, 0, len(txs))
	for _, t := range txs {
		resp = append(resp, dto.BookingItemResponse{
			OrderID:                t.OrderID,
			TransactionStatus:      t.TransactionStatus,
			BookingLifecycleStatus: t.BookingLifecycleStatus,
			StartDate:              t.StartDate,
			EndDate:                t.EndDate,
			TotalPrice:             t.TotalPrice,
			Item: dto.BookingItemInfo{
				ID:    t.ItemID,
				Name:  t.Item.Name,
				Photo: pickFirstPhotoURL(t.Item.PhotoURL),
			},
		})
	}

	return helper.Response{
		Status:   http.StatusOK,
		Messages: "Berhasil mengambil transaksi",
		Data:     resp,
	}
}

func pickFirstPhotoURL(photoURL string) string {
	if photoURL == "" {
		return ""
	}
	var arr []string
	_ = json.Unmarshal([]byte(photoURL), &arr)
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}

func (s *transactionService) SetOngoing(orderID string) helper.Response {
	if err := s.txRepo.UpdateTransactionStatus(orderID, model.TransactionOngoing); err != nil {
		return helper.Response{Status: http.StatusInternalServerError, Messages: "Gagal update status", Error: err.Error()}
	}
	return helper.Response{Status: http.StatusOK, Messages: "Status diubah ke ONGOING"}
}

func (s *transactionService) SetCompleted(orderID string) helper.Response {
	if err := s.txRepo.UpdateTransactionStatus(orderID, model.TransactionCompleted); err != nil {
		return helper.Response{Status: http.StatusInternalServerError, Messages: "Gagal update status", Error: err.Error()}
	}
	return helper.Response{Status: http.StatusOK, Messages: "Status diubah ke COMPLETED"}
}

func (s *transactionService) FinishBooking(orderID string) helper.Response {
	if err := s.txRepo.UpdateBookingLifecycle(orderID, model.BookingFinished); err != nil {
		return helper.Response{Status: http.StatusInternalServerError, Messages: "Gagal menutup booking", Error: err.Error()}
	}
	return helper.Response{Status: http.StatusOK, Messages: "Booking FINISHED"}
}

func (s *transactionService) verifySignature(req dto.MidtransWebhookRequest) bool {
	raw := req.OrderID + req.StatusCode + req.GrossAmount + s.midtransConfig.Serverkey
	hash := sha512.Sum512([]byte(raw))
	return hex.EncodeToString(hash[:]) == req.SignatureKey
}

func mapMidtransStatus(transactionStatus, fraudStatus string) model.PaymentStatus {
	switch transactionStatus {
	case "capture":
		if fraudStatus == "challenge" {
			return model.PaymentPending
		}
		return model.PaymentPaid
	case "settlement":
		return model.PaymentPaid
	case "pending":
		return model.PaymentPending
	case "expire":
		return model.PaymentExpired
	case "deny", "cancel":
		return model.PaymentFailed
	default:
		return model.PaymentPending
	}
}
