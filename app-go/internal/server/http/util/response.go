package util

import (
	"encoding/json"
	"net/http"
)

type BaseResponse[T any] struct {
	Error string `json:"error"`
	Data  T      `json:"data"`
}

// RespondErrorJson is a generic handler to return error json response to api requester
func RespondErrorJson(w http.ResponseWriter, httpstatusCode int, message string) {
	// construct body to json's raw message
	raw, _ := json.Marshal(BaseResponse[struct{}]{
		Error: message,
	})

	// write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpstatusCode)
	w.Write(raw)
}

// RespondJson is a generic handler to return json response to api requester
func RespondJson[T any](w http.ResponseWriter, httpstatusCode int, body T) {
	// construct response
	raw, _ := json.Marshal(BaseResponse[T]{
		Data: body,
	})

	// write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpstatusCode)
	w.Write(raw)
}
