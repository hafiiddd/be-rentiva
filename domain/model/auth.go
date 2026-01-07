package model

import "github.com/golang-jwt/jwt/v5"

type Auth struct {
	Iduser int    `json:"id_user" gorm:"column:id_user;primaryKey;autoIncrement"`
	Username  string `json:"username" gorm:"column:username"`
	jwt.RegisteredClaims
}