package model

import "time"

type ConditionType string

const (
	ConditionBeforeSend   ConditionType = "BEFORE_SEND"
	ConditionAfterReceive ConditionType = "AFTER_RECEIVE"
	ConditionBeforeReturn ConditionType = "BEFORE_RETURN"
	ConditionAfterReturn  ConditionType = "AFTER_RETURN"
)

func (c ConditionType) IsValid() bool {
	switch c {
	case ConditionBeforeSend, ConditionAfterReceive, ConditionBeforeReturn, ConditionAfterReturn:
		return true
	default:
		return false
	}
}

type ItemCondition struct {
	IdCondition int `json:"id_condition" gorm:"column:id_condition;primaryKey;autoIncrement"`

	TransactionID int `json:"transaction_id" gorm:"column:transaction_id;not null"`
	UserID        int `json:"user_id" gorm:"column:user_id;not null"`

	ConditionType ConditionType `json:"condition_type" gorm:"column:condition_type;type:varchar(20)"`
	PhotoURL      string `json:"photo_url" gorm:"column:photo_url;type:text"`
	Note          string `json:"note" gorm:"column:note;type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
func (ItemCondition) TableName() string {
	return "item_conditions"
}
