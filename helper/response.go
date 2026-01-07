package helper

type Response struct {
	Status   int         `json:"status"`
	Messages string      `json:"message"`
	Token    string      `json:"token,omitempty"`
	Error    string      `json:"error,omitempty"`
	Data     interface{} `json:"data,omitempty"`
}
