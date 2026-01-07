package model

import "time"

type Item struct {
	Item_ID     uint      `json:"id_item" gorm:"primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Category    string    `json:"category" gorm:"type:varchar(50)"`
	PricePerDay int       `json:"price_per_day" gorm:"not null"`
	PhotoURL    string    `json:"photo_url" gorm:"type:text"` // simpan JSON array of object keys
	Status      string    `json:"status" gorm:"type:varchar(20);check:status IN ('available','rented','unavailable');default:'available'"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoCreateTime"`

	OwnerID uint `json:"owner_id" gorm:"not null"`
	Owner   User `json:"owner" gorm:"foreignKey:OwnerID;references:Iduser"`

	PhotoKeys []string `json:"photo_keys" gorm:"-"`
}
