package dto

import "back-end/domain/model"

type UploadItemConditionRequest struct {
	ConditionType model.ConditionType `form:"condition_type" json:"condition_type" validate:"required"`
	Note          string              `form:"note" json:"note"`
}
