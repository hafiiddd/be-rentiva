package service

import (
	"back-end/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type MidtransService struct {
	Config *config.MidtransConfig
}

// CreateSnap implements [PaymentGateway].
func (m *MidtransService) CreateSnap(orderID string, amount int64) (string, error) {
	now := time.Now()
	// Midtrans memerlukan time string tanpa offset pecahan detik
	startTime := now.Format("2006-01-02 15:04:05 -0700")
	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     orderID,
			"gross_amount": amount,
		},
		// Batasi masa berlaku pembayaran 10 menit
		"expiry": map[string]interface{}{
			"start_time": startTime,
			"unit":       "minutes",
			"duration":   10,
		},
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		m.Config.BaseUri+"/snap/v1/transactions",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(m.Config.Serverkey, "")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	token, _ := res["token"].(string)
	if token == "" {
		return "", fmt.Errorf("token kosong dari midtrans")
	}

	return token, nil
}

type PaymentGateway interface {
	CreateSnap(orderID string, amount int64) (string, error)
}

func NewMidtransService(cfg *config.MidtransConfig) PaymentGateway {
	return &MidtransService{Config: cfg}
}
