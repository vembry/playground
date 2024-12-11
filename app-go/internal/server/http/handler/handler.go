package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// groupMux groups endpoints with 1 prefix
func groupMux(path string, mux *http.ServeMux) *http.ServeMux {
	group := http.NewServeMux()
	group.Handle(fmt.Sprintf("%s/", path), http.StripPrefix(path, mux)) // cover the entrypoint with middleware
	return group
}

type BaseResponse[T any] struct {
	Error  string `json:"error"`
	Object T      `json:"object"`
}

// respondErrorJson is a generic handler to return error json response to api requester
func respondErrorJson(w http.ResponseWriter, httpstatusCode int, message string) {
	// construct body to json's raw message
	raw, _ := json.Marshal(BaseResponse[struct{}]{
		Error: message,
	})

	// write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpstatusCode)
	w.Write(raw)
}

// respondJson is a generic handler to return json response to api requester
func respondJson[T any](w http.ResponseWriter, httpstatusCode int, body T) {
	// construct response
	raw, _ := json.Marshal(BaseResponse[T]{
		Object: body,
	})

	// write to response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpstatusCode)
	w.Write(raw)
}

func readBody[T any](r *http.Request, target T) error {
	// read payload
	bodyraw, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("unable to read request body. err=%w", err)
	}
	defer r.Body.Close()

	err = json.Unmarshal(bodyraw, &target)
	if err != nil {
		return fmt.Errorf("unable to parse request body. err=%w", err)
	}

	return nil
}
