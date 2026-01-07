package config

import "os"

type MidtransConfig struct {
	// butuh serverkey dan endpoint dev nya midtrans
	Serverkey string
	BaseUri   string
}

func LoadMidtransConfig() MidtransConfig {
	return MidtransConfig{
		Serverkey: os.Getenv("MIDTRANS_SERVER_KEY"),
		BaseUri: "https://app.sandbox.midtrans.com",
	}
}