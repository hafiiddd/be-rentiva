package controller

import (
	"back-end/domain/dto"
	"back-end/domain/model"
	"back-end/domain/service"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type ItemConditionController struct {
	itemConditionService service.ItemConditionService
}

func NewItemConditionController(itemConditionService service.ItemConditionService) *ItemConditionController {
	return &ItemConditionController{itemConditionService: itemConditionService}
}

func (ic *ItemConditionController) UploadItemCondition(c echo.Context) error {
	orderID := c.Param("orderID")

	payload := dto.UploadItemConditionRequest{
		ConditionType: model.ConditionType(c.FormValue("condition_type")),
		Note:          c.FormValue("note"),
	}
	if payload.ConditionType == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "condition_type wajib diisi"})
	}

	fileHeader, err := c.FormFile("photo")
	if err != nil || fileHeader == nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "photo wajib diupload"})
	}
	src, err := fileHeader.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "gagal membuka file"})
	}
	defer src.Close()

	userToken, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "token tidak valid"})
	}
	claims, ok := userToken.Claims.(*model.Auth)
	if !ok {
		return c.JSON(http.StatusUnauthorized, echo.Map{"message": "token tidak valid"})
	}

	res := ic.itemConditionService.UploadItemCondition(
		c.Request().Context(),
		orderID,
		claims.Iduser,
		payload,
		service.UploadConditionFile{
			Filename:    fileHeader.Filename,
			Size:        fileHeader.Size,
			ContentType: fileHeader.Header.Get("Content-Type"),
			Reader:      src,
		},
	)
	return c.JSON(http.StatusOK, res)
}
