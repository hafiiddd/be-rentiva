package controller

import (
	"back-end/domain/dto"
	"back-end/domain/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type PaymentController struct {
	paymentService service.PaymentService
}

func NewPaymentController(paymentService service.PaymentService) *PaymentController {
	return &PaymentController{paymentService: paymentService}
}

func (pc *PaymentController) SaveSnapResult(c echo.Context) error {
	payload := new(dto.SnapResultRequest)
	if err := c.Bind(payload); err != nil {
		return err
	}
	if err := c.Validate(payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "ada bagian yg kosong atau format salah",
		})
	}

	res := pc.paymentService.SaveSnapResult(*payload)
	return c.JSON(http.StatusOK, res)
}
