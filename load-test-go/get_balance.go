package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type response[T any] struct {
	Error  string `json:"error"`
	Object T      `json:"object"`
}

type balance struct {
	Id        string    `json:"id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetBalance(balanceId string) *balance {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf("http://localhost:8080/balance/%s", balanceId),
		nil,
	)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error on http request. error=%v", err)
		return nil
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("get balance '%s' return non 200", balanceId)
		return nil
	}

	rawbody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("error on reading body. error=%v", err)
		return nil
	}

	var out response[balance]
	_ = json.Unmarshal(rawbody, &out)

	return &out.Object
}
