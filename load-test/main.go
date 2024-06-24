package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
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

const AppHost string = "http://host.docker.internal:8080"

var balanceIds []string = []string{
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
}

type traceHandler struct {
	slog.Handler
}

func (h *traceHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}
	return h.Handler.Handle(ctx, r)
}

func main() {

	// exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	// if err != nil {
	// 	log.Fatalf("error initiating trace. err=%v", err)
	// }

	tp := sdktrace.NewTracerProvider(
		// sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("load-test"),
		)),
	)
	otel.SetTracerProvider(tp)

	defer (func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	})()

	tracer := otel.Tracer("load-test-tracer")

	end := time.Now().Add(1 * time.Minute)
	fmt.Println("starting")

	for time.Now().Before(end) {
		for _, balanceId := range balanceIds {
			ctx, span := tracer.Start(context.Background(), "main")
			defer span.End()

			handler := slog.NewJSONHandler(os.Stdout, nil)
			customHandler := &traceHandler{Handler: handler}
			logger := slog.New(customHandler)

			// logger.InfoContext(ctx, "starting for", slog.String("balanceId", balanceId))

			// get balance
			bal, err := GetBalance(ctx, balanceId)
			if err != nil {
				logger.ErrorContext(ctx, "got error on get-balance. err=%v", slog.Any("error", err))
				continue
			}
			if bal == nil {
				logger.ErrorContext(ctx, "get-balance return nil")
				continue
			}

			amount := rand.Float64() * 1000

			// deposit money when needed
			if bal.Amount < amount {
				depositAmount := 100 + rand.Float64()*10000
				err := Deposit(ctx, balanceId, depositAmount)
				if err != nil {
					logger.ErrorContext(ctx, "got error on deposit. err=%v", slog.Any("error", err))
					continue
				}
			}

			// exec withdrawal
			err = Withdraw(ctx, balanceId, amount)
			if err != nil {
				logger.ErrorContext(ctx, "got error on withdrawal. err=%v", slog.Any("error", err))
				continue
			}
		}
	}

	fmt.Println("finished")
}

func Deposit(ctx context.Context, balanceId string, amount float64) error {
	payloadRaw, _ := json.Marshal(map[string]interface{}{
		"amount": amount,
	})

	req, _ := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/balance/%s/deposit", AppHost, balanceId),
		bytes.NewBuffer(payloadRaw),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error on http request. error=%w", err)
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

func GetBalance(ctx context.Context, balanceId string) (*balance, error) {
	req, _ := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("%s/balance/%s", AppHost, balanceId),
		nil,
	)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error on http request. error=%w", err)
	}

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

func Withdraw(ctx context.Context, balanceId string, amount float64) error {
	payloadRaw, _ := json.Marshal(map[string]interface{}{
		"amount": amount,
	})

	req, _ := http.NewRequestWithContext(
		ctx,
		"POST",
		fmt.Sprintf("%s/balance/%s/withdraw", AppHost, balanceId),
		bytes.NewBuffer(payloadRaw),
	)
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("error on http request. error=%w", err)
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
