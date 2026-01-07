package service

import (
	"back-end/domain/model"
	"back-end/domain/repository"
	"back-end/helper"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepository repository.UserRepository
}

// DisplayAccount implements UserService.
func (u *userService) DisplayAccount(CurrentUserID int) helper.Response {
	items, err := u.userRepository.DisplayAccount(CurrentUserID)
	if err != nil {
		return helper.Response{
			Status:   http.StatusInternalServerError,
			Messages: "Gagal mengambil data barang di database",
			Error:    err.Error(),
		}
	}

	return helper.Response{
		Status:   http.StatusCreated,
		Messages: "Berhasil mengambil data account",
		Data:     items,
	}
}

// Login implements UserService.
func (u *userService) Login(username string, password string) helper.Response {
	//cari akun udah ada apa belom
	data, err := u.userRepository.FindUser(username)
	if err != nil {
		//user ga ditemukan
		return helper.Response{
			Status:   400,
			Messages: "username tidak ditemukan",
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(password))
	if err != nil {
		return helper.Response{
			Status:   400,
			Messages: "Password tidak sesuai",
		}
	}
	claims := &model.Auth{
		Iduser:   data.Iduser,
		Username: data.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	secretFromEnv := strings.TrimSpace(os.Getenv("JWT_SECRET_KEY"))
	jwtSecretKey := []byte(secretFromEnv)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signToken, err := token.SignedString(jwtSecretKey)

	if err != nil {
		return helper.Response{
			Status:   400,
			Messages: "Gagal login",
		}
	}
	return helper.Response{
		Status:   201,
		Messages: "berhasil login",
		Token:    signToken,
		Data:     data,
	}
}

// Register implements UserService.
func (u *userService) Register(newUser model.User) helper.Response {
	//cek dulu apakah ada user lain
	_, err := u.userRepository.FindUser(newUser.Email)
	if err == nil {
		return helper.Response{
			Status:   400,
			Messages: "email sudah terdaftar",
		}
	}
	//hash pass
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		return helper.Response{
			Status:   400,
			Messages: "gagal hash password",
		}
	}
	//create akun
	data, err := u.userRepository.Register(model.User{
		Email:    newUser.Email,
		Password: string(hash),
		Username: newUser.Username,
		FullName: newUser.FullName,
		Nik:      newUser.Nik,
		Ttl:      newUser.Ttl,
		Address:  newUser.Address,
	})
	//response
	if err != nil {
		return helper.Response{
			Status:   400,
			Messages: "gagal register",
		}
	}
	return helper.Response{
		Status:   200,
		Messages: "berhasil register",
		Data:     data,
	}
}

type UserService interface {
	Login(email, password string) helper.Response
	Register(newUser model.User) helper.Response
	DisplayAccount(CurrenUserID int) helper.Response
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userService{userRepository: userRepository}
}
