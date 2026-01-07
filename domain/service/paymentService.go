package service

import (
	"back-end/domain/dto"
	"back-end/domain/repository"
	"back-end/helper"
	"net/http"
	"time"
)

type PaymentService interface {
	SaveSnapResult(req dto.SnapResultRequest) helper.Response
}

type paymentService struct {
	paymentRepo repository.PaymentRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository) PaymentService {
	return &paymentService{paymentRepo: paymentRepo}
}

func (s *paymentService) SaveSnapResult(req dto.SnapResultRequest) helper.Response {
	expiry := time.Now().Add(10 * time.Minute)

	updates := map[string]interface{}{
		"payment_type":   req.PaymentType,
		"payment_method": req.PaymentMethod,
		"va_number":      req.VANumber,
		"expiry_time":    expiry,
	}

	payment, err := s.paymentRepo.UpdateSnapInfoByOrderID(req.OrderID, updates)
	if err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal menyimpan hasil snap",
			Error:    err.Error(),
		}
	}

	return helper.Response{
		Status:   http.StatusOK,
		Messages: "Berhasil menyimpan instrumen pembayaran",
		Data:     payment,
	}
}
