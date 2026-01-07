package model

import "time"

type Review struct {
	IDReview     int       `json:"id_review" gorm:"primaryKey;autoIncrement"`
	TransactionID int      `json:"transaction_id" gorm:"not null"`
	ReviewerID    int      `json:"reviewer_id" gorm:"not null"`
	RevieweeID    int      `json:"reviewee_id" gorm:"not null"`
	Rating        int      `json:"rating" gorm:"not null"`
	Comment       string   `json:"comment" gorm:"type:text"`
	RoleContext   string   `json:"role_context" gorm:"type:varchar(10)"` // RENTER / OWNER (opsional audit)
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Review) TableName() string {
	return "reviews"
}
