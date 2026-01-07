package main

import (
	"back-end/config"
	"back-end/domain/controller"
	"back-end/domain/model"
	"back-end/domain/repository"
	"back-end/domain/service"
	"back-end/router"
	validatorpkg "back-end/validator"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {

	errs := godotenv.Load()
	if errs != nil {
		log.Println("Warning: Error loading .env file, using default environment variables")
	}
	db := config.InitDb()
	err := db.AutoMigrate(&model.User{}, &model.Item{}, &model.Transaction{}, &model.Payment{}, &model.Review{}, &model.ItemCondition{})
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}

	//dependecies injection
	authRepo := repository.NewUserRepository(db)
	authService := service.NewUserService(authRepo)
	authCtrl := controller.NewUserController(authService)

	itemRepo := repository.NewItemRepository(db)
	itemService := service.NewItemService(itemRepo)

	minioCfg := config.InitMinio()
	storageSvc := service.NewFileService(minioCfg)
	itemCtrl := controller.NewItemController(itemService, storageSvc)

	txRepo := repository.NewTransactionRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	reviewRepo := repository.NewReviewRepository(db)
	itemConditionRepo := repository.NewItemConditionRepository(db)

	midCfg := config.LoadMidtransConfig()
	paymentGateway := service.NewMidtransService(&midCfg)
	txService := service.NewTransactionService(txRepo, paymentRepo, paymentGateway, midCfg)
	txCtrl := controller.NewTransactionController(txService)
	webhookCtrl := controller.NewWebhookController(txService)
	paymentService := service.NewPaymentService(paymentRepo)
	paymentCtrl := controller.NewPaymentController(paymentService)
	reviewService := service.NewReviewService(reviewRepo, txRepo)
	reviewCtrl := controller.NewReviewController(reviewService)
	itemConditionService := service.NewItemConditionService(txRepo, itemConditionRepo, storageSvc)
	itemConditionCtrl := controller.NewItemConditionController(itemConditionService)
	e := echo.New()

	e.Validator = validatorpkg.New()
	router.SetUp(e, authCtrl, itemCtrl, txCtrl, webhookCtrl, paymentCtrl, reviewCtrl, itemConditionCtrl)
	port := ":8080"
	log.Printf("Server berjalann di port %s", port)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	if err := e.Start(port); err != nil && err != http.ErrServerClosed {
		e.Logger.Fatal("Gagal menjalankan server: ", err)
	}

}
