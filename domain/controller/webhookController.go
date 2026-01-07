package controller

import (
	"back-end/domain/dto"
	"back-end/domain/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type WebhookController struct {
	transactionService service.TransactionService
}

func NewWebhookController(transactionService service.TransactionService) *WebhookController {
	return &WebhookController{transactionService: transactionService}
}

func (wc *WebhookController) MidtransCallback(c echo.Context) error {
	payload := new(dto.MidtransWebhookRequest)
	log.Println("🔥 MIDTRANS WEBHOOK MASUK")
	if err := c.Bind(payload); err != nil {
		return err
	}

	res := wc.transactionService.HandleWebhook(*payload)
	return c.JSON(http.StatusOK, res)
}
