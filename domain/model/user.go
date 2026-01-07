package model

import "time"

type User struct {
	Iduser         int    `json:"id_user" gorm:"column:id_user;primaryKey;autoIncrement"`
	Username       string `json:"username" gorm:"column:username;type:varchar(50);unique;not null"`
	Email          string `json:"email" gorm:"column:email"`
	Password       string `json:"password" gorm:"column:password"`
	FullName       string `json:"full_name" gorm:"column:full_name"`
	Nik            string `json:"nik" gorm:"column:nik;type:varchar(20)"`
	Ttl            string `json:"ttl" gorm:"column:ttl;type:varchar(100)"`
	Address        string `json:"address" gorm:"column:address;type:text"`
	KtpImageUrl    string `json:"ktp_image_url" gorm:"column:ktp_image_url;type:text"`
	SelfieImageUrl string `json:"selfie_image_url" gorm:"column:selfie_image_url;type:text"`

	VerificationStatus     string  `json:"verification_status" gorm:"column:verification_status;type:varchar(20);default:'PENDING'"`
	TrustScore             float64 `json:"trust_score" gorm:"column:trust_score;default:70.0"`
	TotalTransactions      int     `json:"total_transactions" gorm:"column:total_transactions;default:0"`
	SuccessfulTransactions int     `json:"successful_transactions" gorm:"column:successful_transactions;default:0"`

	Created_at time.Time
	Updated_at time.Time

	Items []Item `json:"items" gorm:"foreignKey:OwnerID;references:Iduser"`
}

func (User) TableName() string {
	return "users"
}

// -- Data dari Langkah 4 (Validasi Data Manual)
// full_name VARCHAR(255),  -- Dari 'finalNama'
// nik VARCHAR(20),         -- Dari 'finalNik'
// ttl VARCHAR(100),        -- Dari 'finalTtl'
// address TEXT,            -- Dari 'finalAlamat'

// -- Path/URL ke file (disimpan di GCS/S3/Minio, BUKAN di database)
// ktp_image_url TEXT,      -- Dari 'ktpImage' (Langkah 2)
// selfie_image_url TEXT,   -- Dari 'selfieImage' (Langkah 3)

// -- -----------------------------------------
// -- KOLOM PENTING UNTUK TRUST SYSTEM
// -- -----------------------------------------

// -- Status untuk Verifikasi Manual oleh Admin
// -- ( 'PENDING', 'VERIFIED', 'REJECTED' )
// verification_status VARCHAR(20) DEFAULT 'PENDING',

// -- Kolom untuk 'Trust Score' (Future Improvement)
// trust_score DOUBLE PRECISION DEFAULT 70.0, -- Skor awal
// total_transactions INT DEFAULT 0,
// successful_transactions INT DEFAULT 0,

// -- Timestamps
// created_at TIMESTAMPTZ DEFAULT NOW(),
// updated_at TIMESTAMPTZ DEFAULT NOW()
