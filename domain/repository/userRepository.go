package repository

import (
	"back-end/domain/model"


	"gorm.io/gorm"
)

type dbUser struct {
	conn *gorm.DB
}

// DisplayAccount implements UserRepository.
func (d *dbUser) DisplayAccount(CurrentUserId int) ([]model.User, error) {
	var data []model.User
	err := d.conn.Where("id_user != ?", CurrentUserId).First(&data).Error
	if err != nil {
		return nil, err
	}
	return data, err
}

// FindUser implements UserRepository.
func (d *dbUser) FindUser(username string) (model.User, error) {
	var data model.User
	err := d.conn.Where("username = ?", username).First(&data).Error
	if err != nil {
		return model.User{}, err
	}
	return data, nil
}

// Register implements UserRepository.
func (d *dbUser) Register(newUser model.User) (model.User, error) {
	
	err := d.conn.Create(&newUser).Error
	if err != nil {
		return model.User{}, err
	}
	return newUser, nil
}

type UserRepository interface {
	Register(newUser model.User) (model.User, error)
	FindUser(username string) (model.User, error)
	DisplayAccount(CurrentUserId int) ([]model.User, error)
}

func NewUserRepository(conn *gorm.DB) UserRepository {
	return &dbUser{conn: conn}
}
