package service

import (
	"back-end/domain/dto"
	"back-end/domain/model"
	"back-end/domain/repository"
	"back-end/helper"
	"net/http"
)

type ReviewService interface {
	CreateReview(req dto.CreateReviewRequest, reviewerID int) helper.Response
	GetByTransaction(txID int) helper.Response
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
	txRepo     repository.TransactionRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository, txRepo repository.TransactionRepository) ReviewService {
	return &reviewService{reviewRepo: reviewRepo, txRepo: txRepo}
}

func (s *reviewService) CreateReview(req dto.CreateReviewRequest, reviewerID int) helper.Response {
	tx, err := s.txRepo.FindByID(req.TransactionID)
	if err != nil {
		return helper.Response{Status: http.StatusBadRequest, Messages: "Transaksi tidak ditemukan", Error: err.Error()}
	}

	// Validasi status transaksi
	if tx.TransactionStatus != model.TransactionCompleted || tx.BookingLifecycleStatus != model.BookingInProgress {
		return helper.Response{Status: http.StatusBadRequest, Messages: "Review hanya bisa untuk transaksi COMPLETED & booking IN_PROGRESS"}
	}

	// Tentukan reviewee
	var revieweeID int
	roleContext := ""
	ownerID := int(tx.Item.OwnerID)
	switch reviewerID {
	case tx.RenterID:
		revieweeID = ownerID
		roleContext = "RENTER"
	case ownerID:
		revieweeID = tx.RenterID
		roleContext = "OWNER"
	default:
		return helper.Response{Status: http.StatusForbidden, Messages: "Anda bukan pihak transaksi ini"}
	}

	if revieweeID == reviewerID {
		return helper.Response{Status: http.StatusBadRequest, Messages: "Tidak dapat review diri sendiri"}
	}

	// Cek sudah pernah review?
	hasReviewed, err := s.reviewRepo.HasUserReviewed(req.TransactionID, reviewerID)
	if err != nil {
		return helper.Response{Status: http.StatusInternalServerError, Messages: "Gagal cek review", Error: err.Error()}
	}
	if hasReviewed {
		return helper.Response{Status: http.StatusBadRequest, Messages: "Anda sudah memberi review untuk transaksi ini"}
	}

	if req.Rating < 1 || req.Rating > 5 {
		return helper.Response{Status: http.StatusBadRequest, Messages: "Rating harus 1-5"}
	}

	review := model.Review{
		TransactionID: req.TransactionID,
		ReviewerID:    reviewerID,
		RevieweeID:    revieweeID,
		Rating:        req.Rating,
		Comment:       req.Comment,
		RoleContext:   roleContext,
	}
	review, err = s.reviewRepo.Create(review)
	if err != nil {
		return helper.Response{Status: http.StatusInternalServerError, Messages: "Gagal menyimpan review", Error: err.Error()}
	}

	// Cek kedua review sudah ada
	if err := s.finishIfBothReviewed(req.TransactionID, tx.OrderID); err != nil {
		return helper.Response{Status: http.StatusInternalServerError, Messages: "Gagal update status booking", Error: err.Error()}
	}

	return helper.Response{Status: http.StatusCreated, Messages: "Review berhasil", Data: review}
}

func (s *reviewService) GetByTransaction(txID int) helper.Response {
	reviews, err := s.reviewRepo.FindByTransaction(txID)
	if err != nil {
		return helper.Response{Status: http.StatusInternalServerError, Messages: "Gagal mengambil review", Error: err.Error()}
	}
	return helper.Response{Status: http.StatusOK, Messages: "Berhasil mengambil review", Data: reviews}
}

func (s *reviewService) finishIfBothReviewed(txID int, orderID string) error {
	count, err := s.reviewRepo.CountByTransaction(txID)
	if err != nil {
		return err
	}
	// butuh 2 review: renter -> owner, owner -> renter
	if count < 2 {
		return nil
	}
	// update booking lifecycle ke FINISHED
	return s.txRepo.UpdateBookingLifecycle(orderID, model.BookingFinished)
}
