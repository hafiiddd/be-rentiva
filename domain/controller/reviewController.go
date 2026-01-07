package controller

import (
	"back-end/domain/dto"
	"back-end/domain/model"
	"back-end/domain/service"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ReviewController struct {
	reviewService service.ReviewService
}

func NewReviewController(reviewService service.ReviewService) *ReviewController {
	return &ReviewController{reviewService: reviewService}
}

func (rc *ReviewController) CreateReview(c echo.Context) error {
	payload := new(dto.CreateReviewRequest)
	if err := c.Bind(payload); err != nil {
		return err
	}
	if err := c.Validate(payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "payload tidak valid"})
	}

	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "token tidak valid"})
	}
	claims, ok := userToken.Claims.(*model.Auth)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "token tidak valid"})
	}

	res := rc.reviewService.CreateReview(*payload, claims.Iduser)
	return c.JSON(http.StatusOK, res)
}

func (rc *ReviewController) GetByTransaction(c echo.Context) error {
	txIDStr := c.Param("transaction_id")
	txID, err := strconv.Atoi(txIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "transaction_id tidak valid"})
	}
	res := rc.reviewService.GetByTransaction(txID)
	return c.JSON(http.StatusOK, res)
}
