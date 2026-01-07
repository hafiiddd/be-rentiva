package repository

import (
	"back-end/domain/model"

	"gorm.io/gorm"
)

type ReviewRepository interface {
	Create(review model.Review) (model.Review, error)
	HasUserReviewed(txID int, userID int) (bool, error)
	CountByTransaction(txID int) (int64, error)
	FindByTransaction(txID int) ([]model.Review, error)
}

type dbReview struct {
	conn *gorm.DB
}

func NewReviewRepository(conn *gorm.DB) ReviewRepository {
	return &dbReview{conn: conn}
}

func (r *dbReview) Create(review model.Review) (model.Review, error) {
	if err := r.conn.Create(&review).Error; err != nil {
		return model.Review{}, err
	}
	return review, nil
}

func (r *dbReview) HasUserReviewed(txID int, userID int) (bool, error) {
	var count int64
	if err := r.conn.Model(&model.Review{}).
		Where("transaction_id = ? AND reviewer_id = ?", txID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *dbReview) CountByTransaction(txID int) (int64, error) {
	var count int64
	if err := r.conn.Model(&model.Review{}).
		Where("transaction_id = ?", txID).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *dbReview) FindByTransaction(txID int) ([]model.Review, error) {
	var reviews []model.Review
	if err := r.conn.Where("transaction_id = ?", txID).Find(&reviews).Error; err != nil {
		return nil, err
	}
	return reviews, nil
}
