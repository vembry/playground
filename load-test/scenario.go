package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

type parameter struct {
	BalanceIds []string
}

func initiateParameter() parameter {
	return parameter{
		BalanceIds: []string{
			"2TeSprhp2cN6nEIcayZsjjvnlsK",
			"2TeSppzLaJaxldlzMOkqYO37vqw",
			"2TeSprCLB0tF6HsJT9eCb1IBht0",
			"2TeSps1ECSPx2IRrWMKEd6oHvSJ",
			"2TeSpqL5Cq3t0rNJ4RDcR4ky029",
			"2TeSpo6Okj6BedA62PUO3nRYADm",
			"2TeSprwUG0oLfgGbIIObipgf5be",
			"2TeSpqi0vPnYFf82cdBCE4Gaxjx",
			"2TeSppFijPxu2bIKpRLhLy3ye4j",
			"2TeSpnNHCn0Yq0vJddsv0dcSviy",
			"2TeSpsBbTYE5ZBe3zpO1HDRf44o",
			// "2TeSprmZUc7HGpqlm7bUi6RB9lu",
			// "2TeSpqKkgW9pz15PaqR8ZIS0Wjh",
			// "2TeSpqfWmsmNLffPB4fhErgYRiC",
			// "2TeSpnH3d9Hr9zPyNpoBrzdd8fx",
			// "2TeSpsz6jeis3vJHCoSfFrnGPVl",
			// "2TeSptEF7KEwZ9cdlowOH4HSmbR",
			// "2TeSpsvjaufyDq78BwcFFZ6ibfi",
			// "2TeSpsFLvMYhPU3CpC21wgyOcS5",
			// "2TeSpmhKLpqgNXTNOvh3JPLkbjy",
			// "2TWlPQ2AhstX9PtJ5UTOE6xQ7Ga",
			// "2TWlPVmPjhonQ2DpOFt09O990th",
			// "2TWlPdYWFbP3iXPMIKVRmdZ3ozC",
			// "2TWlPkjd1YcRDCBDQk11nygVDpe",
			// "2TWlPw3U64nhEtzeazP5ELd7q4c",
			// "2TWlQ2JWBetv9MUkFsLPd4zhLa4",
			// "2TWlQ83iEHKIOtvMCfSnBwM2sEB",
			// "2TWlQJsEhRIW8XpeGGz2u75phWN",
			// "2TWlQRnYNr5ViFi0wLTatYImXz5",
			// "2TWlQYYZVIC566T3XFVldQckPsB",
			// "2TWlQcFLvqE67A2qbZpWAdSJiZL",
			// "2TWlQjcMepYBrRJRRzwgXWLA0gX",
			// "2TWlQvGoGcHxUa4iod28bMsfW7e",
			// "2TWlR5WFe7VUVMSxhgpEmfqlAAX",
			// "2TWlR8ifwogFmktTwET0Eb2s4PE",
			// "2TWlRFfCwsZo903aTO7xSRCYQIU",
			// "2TWlRRtrVofuy7C1ZzcWIICVEME",
			// "2TWlRZBGQlbQe34dVNU3GhQKshe",
			// "2TWlRedfruxmFYvYLJur7oGesXY",
			// "2TWlRnR3C0hopS7NkcwyjIOq5Kd",
		},
	}
}

func GetBalance(ctx context.Context, balanceId string) (*balance, error) {
	res, err := otelhttp.Get(
		ctx,
		fmt.Sprintf("%s/balance/%s", AppHost, balanceId),
	)
	if err != nil {
		return nil, fmt.Errorf("error on http request. error=%w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get balance '%s' return non 200", balanceId)
	}

	rawbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading body. error=%w", err)
	}

	var out response[balance]
	_ = json.Unmarshal(rawbody, &out)

	return &out.Object, nil
}

func Deposit(ctx context.Context, balanceId string, amount float64) error {
	payloadRaw, _ := json.Marshal(map[string]interface{}{
		"amount": amount,
	})

	res, err := otelhttp.Post(
		ctx,
		fmt.Sprintf("%s/balance/%s/deposit", AppHost, balanceId),
		"application/json",
		bytes.NewBuffer(payloadRaw),
	)
	if err != nil {
		return fmt.Errorf("error on http request. error=%w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("depositing balance '%s' return non 200", balanceId)
	}

	return nil
}

func Withdraw(ctx context.Context, balanceId string, amount float64) error {
	payloadRaw, _ := json.Marshal(map[string]interface{}{
		"amount": amount,
	})

	res, err := otelhttp.Post(
		ctx,
		fmt.Sprintf("%s/balance/%s/withdraw", AppHost, balanceId),
		"application/json",
		bytes.NewBuffer(payloadRaw),
	)
	if err != nil {
		return fmt.Errorf("error on http request. error=%w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("depositing balance '%s' return non 200", balanceId)
	}

	return nil
}
