package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func Withdraw(balanceId string, amount float64) error {
	payloadRaw, _ := json.Marshal(map[string]interface{}{
		"amount": amount,
	})

	req, _ := http.NewRequest(
		"POST",
		fmt.Sprintf("http://localhost:8080/balance/%s/withdraw", balanceId),
		bytes.NewBuffer(payloadRaw),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error on http request. error=%v", err)
		return err
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("depositing balance '%s' return non 200", balanceId)
		return err
	}

	// rawbody, err := io.ReadAll(res.Body)
	// if err != nil {
	// 	log.Printf("error on reading body. error=%v", err)
	// 	return err
	// }

	// var out response[balance]
	// _ = json.Unmarshal(rawbody, &out)

	return nil
}
