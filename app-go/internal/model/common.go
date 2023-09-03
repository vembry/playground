package model

// CommonResponse contain generic field for response
type CommonResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}
