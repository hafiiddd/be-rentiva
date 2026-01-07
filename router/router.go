package router

import (
	"back-end/domain/controller"
	"back-end/middleware"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SetUp(
	e *echo.Echo,
	auth *controller.UserController,
	item *controller.ItemController,
	transaction *controller.TransactionController,
	webhook *controller.WebhookController,
	payment *controller.PaymentController,
	review *controller.ReviewController,
	itemCondition *controller.ItemConditionController,
) {
	e.POST("api/v1/register", auth.Register)
	e.POST("api/v1/login", auth.Login)

	jwtMiddleware := echojwt.WithConfig(middleware.JWTMiddlewareConfig())

	api := e.Group("/api/v1")

	api.Use(middleware.SkipOptions(jwtMiddleware))

	api.GET("/display/:id", auth.DisplayAccount)
	api.GET("/item/:id", item.GetItemByUserID)
	api.GET("/item/detail/:idItem", item.GetItemByID)
	api.POST("/item", item.CreateItem)
	api.PUT("/item/:idItem", item.UpdateItem)
	api.DELETE("/item/:idItem", item.DeleteItem)

	api.POST("/transaction", transaction.CreateTransaction)
	api.POST("/payments/snap-result", payment.SaveSnapResult)

	
	api.GET("/transaction/paid", transaction.ListPaidTransactions)
	api.POST("/transaction/:orderID/item-condition", itemCondition.UploadItemCondition)
	api.PUT("/transaction/:orderID/ongoing", transaction.SetOngoing)
	api.PUT("/transaction/:orderID/completed", transaction.SetCompleted)
	api.PUT("/transaction/:orderID/finish", transaction.FinishBooking)
	api.POST("/reviews", review.CreateReview)
	api.GET("/reviews/transaction/:transaction_id", review.GetByTransaction)

	e.POST("/api/v1/payment/webhook/midtrans", webhook.MidtransCallback)
}
