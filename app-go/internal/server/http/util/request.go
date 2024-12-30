package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ReadBody[T any](r *http.Request, target T) error {
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
