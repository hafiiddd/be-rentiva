package controller

import (
	"back-end/domain/dto"
	"back-end/domain/model"
	"back-end/domain/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// ItemHandler menampung dependensi service
type ItemController struct {
	itemService service.ItemService
	storage     service.StorageService
}

func (h *ItemController) CreateItem(c echo.Context) error {
	payload := new(dto.ItemDTO)
	if err := c.Bind(payload); err != nil {
		return err
	}
	if err := c.Validate(payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "ada bagian yg kosong",
		})
	}

	photoKeys, err := h.collectPhotos(c)
	if err != nil {
		return err
	}
	photoJSON := marshalPhotoURLs(photoKeys)

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

	result := h.itemService.CreateItem(model.Item{
		Name:        payload.Name,
		Description: payload.Description,
		Category:    payload.Category,
		PricePerDay: payload.PricePerDay,
		PhotoURL:    photoJSON,
		Status:      payload.Status,
		OwnerID:     uint(claims.Iduser),
	})

	return c.JSON(http.StatusOK, result)
}

func (h *ItemController) UpdateItem(c echo.Context) error {
	idStr := c.Param("idItem")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid item id",
		})
	}

	payload := new(dto.ItemDTO)
	if err := c.Bind(payload); err != nil {
		return err
	}
	if err := c.Validate(payload); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "ada bagian yg kosong",
		})
	}

	photoKeys, err := h.collectPhotos(c)
	if err != nil {
		return err
	}
	photoJSON := marshalPhotoURLs(photoKeys)

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

	result := h.itemService.UpdateItem(id, model.Item{
		Name:        payload.Name,
		Description: payload.Description,
		Category:    payload.Category,
		PricePerDay: payload.PricePerDay,
		PhotoURL:    photoJSON,
		Status:      payload.Status,
		OwnerID:     uint(claims.Iduser),
	})

	return c.JSON(http.StatusOK, result)
}

func (h *ItemController) DeleteItem(c echo.Context) error {
	idStr := c.Param("idItem")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid item id",
		})
	}

	result := h.itemService.DeleteItem(id)
	return c.JSON(http.StatusOK, result)
}

// GetItemByID (Fungsi untuk Halaman Detail Barang)
func (h *ItemController) GetItemByUserID(c echo.Context) error {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid item id",
		})
	}

	// 3. Panggil service
	result := h.itemService.GetItemsByUserID(id)

	// Auto-generate URL untuk list item
	if items, ok := result.Data.([]model.Item); ok {
		resp := make([]echo.Map, 0, len(items))
		for _, it := range items {
			resp = append(resp, h.withPhotos(it))
		}
		result.Data = resp
	}

	return c.JSON(http.StatusOK, result)
}

func (h *ItemController) GetItemByID(c echo.Context) error {
	idStr := c.Param("idItem")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid item id",
		})
	}

	// 3. Panggil service
	result := h.itemService.GetItemByid(id)

	if item, ok := result.Data.(model.Item); ok {
		result.Data = h.withPhotos(item)
	}

	return c.JSON(http.StatusOK, result)
}

// NewItemController membuat instance Controller baru
func NewItemController(itemService service.ItemService, storage service.StorageService) *ItemController {
	return &ItemController{itemService: itemService, storage: storage}
}

// collectPhotos mengumpulkan object key dari upload multipart (photos/photo).
func (h *ItemController) collectPhotos(c echo.Context) ([]string, error) {
	keys := make([]string, 0)

	// Multiple files key "photos"
	if form, _ := c.MultipartForm(); form != nil {
		if files, ok := form.File["photos"]; ok {
			for _, fh := range files {
				src, err := fh.Open()
				if err != nil {
					return nil, c.JSON(http.StatusBadRequest, echo.Map{"message": "gagal membuka file"})
				}
				defer src.Close()

				objectName := service.UniqueObjectName("items", fh.Filename)
				key, err := h.storage.Upload(c.Request().Context(), objectName, src, fh.Size, fh.Header.Get("Content-Type"))
				if err != nil {
					return nil, c.JSON(http.StatusInternalServerError, echo.Map{"message": "gagal upload ke storage", "error": err.Error()})
				}
				keys = append(keys, key)
			}
		}
	}

	// Fallback single file key "photo" atau "photo_url"
	if len(keys) == 0 {
		if fileHeader, err := c.FormFile("photo"); err == nil && fileHeader != nil {
			src, err := fileHeader.Open()
			if err != nil {
				return nil, c.JSON(http.StatusBadRequest, echo.Map{"message": "gagal membuka file"})
			}
			defer src.Close()

			objectName := service.UniqueObjectName("items", fileHeader.Filename)
			key, err := h.storage.Upload(c.Request().Context(), objectName, src, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
			if err != nil {
				return nil, c.JSON(http.StatusInternalServerError, echo.Map{"message": "gagal upload ke storage", "error": err.Error()})
			}
			keys = append(keys, key)
		}
	}

	return keys, nil
}

func marshalPhotoURLs(urls []string) string {
	if len(urls) == 0 {
		return ""
	}
	if b, err := json.Marshal(urls); err == nil {
		return string(b)
	}
	return ""
}

func parsePhotoKeys(photoURL string, keys []string) []string {
	if len(keys) > 0 {
		return keys
	}
	if photoURL == "" {
		return []string{}
	}
	var arr []string
	_ = json.Unmarshal([]byte(photoURL), &arr)
	return arr
}

func (h *ItemController) buildURLs(keys []string) []string {
	out := make([]string, 0, len(keys))
	for _, k := range keys {
		if k == "" {
			continue
		}
		if urlStr, err := h.storage.BuildPublicURL(k); err == nil {
			out = append(out, urlStr)
			continue
		}
		if urlStr, err := h.storage.GetPresignedURL(k); err == nil {
			out = append(out, urlStr)
		}
	}
	return out
}

func (h *ItemController) withPhotos(item model.Item) echo.Map {
	keys := parsePhotoKeys(item.PhotoURL, item.PhotoKeys)
	urls := h.buildURLs(keys)
	item.PhotoURL = ""
	item.PhotoKeys = keys
	return echo.Map{
		"item":   item,
		"photos": urls,
	}
}
