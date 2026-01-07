package controller

import (
	"back-end/domain/dto"
	"back-end/domain/model"
	"back-end/domain/service"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TransactionController struct {
	transactionService service.TransactionService
}

func NewTransactionController(transactionService service.TransactionService) *TransactionController {
	return &TransactionController{transactionService: transactionService}
}

func (tc *TransactionController) CreateTransaction(c echo.Context) error {
	payload := new(dto.CreateTransactionRequest)
	if err := c.Bind(payload); err != nil {
		return err
	}
	if err := c.Validate(payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "ada bagian yg kosong atau format tanggal salah",
		})
	}

	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "token tidak valid",
		})
	}
	claims, ok := userToken.Claims.(*model.Auth)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "token tidak valid",
		})
	}

	res := tc.transactionService.CreateTransaction(*payload, claims.Iduser)
	return c.JSON(http.StatusOK, res)
}

func (tc *TransactionController) ListPaidTransactions(c echo.Context) error {
	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "token tidak valid",
		})
	}
	claims, ok := userToken.Claims.(*model.Auth)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"message": "token tidak valid",
		})
	}

	res := tc.transactionService.ListPaidTransactions(claims.Iduser)
	return c.JSON(http.StatusOK, res)
}

// opsional: handler untuk update status sewa
func (tc *TransactionController) SetOngoing(c echo.Context) error {
	orderID := c.Param("orderID")
	res := tc.transactionService.SetOngoing(orderID)
	return c.JSON(http.StatusOK, res)
}

func (tc *TransactionController) SetCompleted(c echo.Context) error {
	orderID := c.Param("orderID")
	res := tc.transactionService.SetCompleted(orderID)
	return c.JSON(http.StatusOK, res)
}

func (tc *TransactionController) FinishBooking(c echo.Context) error {
	orderID := c.Param("orderID")
	res := tc.transactionService.FinishBooking(orderID)
	return c.JSON(http.StatusOK, res)
}
