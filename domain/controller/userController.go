package controller

import (
	"back-end/domain/dto"
	"back-end/domain/model"
	"back-end/domain/service"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	userService service.UserService
}

func (uc *UserController) Register(c echo.Context) error {
	payload := new(dto.ResisterDTO)
	if err := c.Bind(payload); err != nil {
		return err
	}
	 if err := c.Validate(payload); err != nil { 
        return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "ada bagian yg kosong",
		})
    }
	fmt.Print(payload)
	res := uc.userService.Register(model.User{
		Email:    payload.Email,
		Password: payload.Password,
		Username: payload.Username,
		FullName: payload.FullName,
		Nik:      payload.Nik,
		Ttl:      payload.Ttl,
		Address:  payload.Address,
	})
	return c.JSON(http.StatusOK, res)
}
func (uc *UserController) Login(c echo.Context) error {
	payload := new(dto.AuthRequest)
	print(payload)
	if err := c.Bind(payload); err != nil {
		return err
	}
	 if err := c.Validate(payload); err != nil { 
        return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "ada bagian yg kosong",
		})
    }
	res := uc.userService.Login(payload.Username, payload.Password)
	return c.JSON(http.StatusOK, res)
}
func (uc *UserController) DisplayAccount(c echo.Context) error{
	idStr := c.Param("id")
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid item id",
		})
	}

	// 3. Panggil service
	result := uc.userService.DisplayAccount(id)
	return c.JSON(http.StatusOK, result)
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}
