package dto

type CreateReviewRequest struct {
	TransactionID int    `json:"transaction_id" validate:"required"`
	Rating        int    `json:"rating" validate:"required,min=1,max=5"`
	Comment       string `json:"comment" validate:"required"`
}
