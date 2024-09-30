package main

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"sdk/tester"
)

func main() {

	// setup parameters
	shutdownHandler := newTelemetry()
	defer shutdownHandler()

	// setup logger
	logger := newLogger()

	// setup load tester
	t := tester.New(
		tester.Config{
			Type:                  tester.LoadType_Count,
			Logger:                logger,
			ConcurrentWorkerCount: 10,
			Counter:               6000,
		},
	)

	// setup parameter for test
	params := initiateParameter()

	// run load tester
	t.Do(func(ctx context.Context, logger *slog.Logger) {
		// choose balance id
		i := randRange(0, len(params.BalanceIds)-1)

		var (
			balanceId = params.BalanceIds[i]
			amount    = rand.Float64() * 1000
			err       error
		)

		// get balance
		bal, err := GetBalance(ctx, balanceId)
		if err != nil {
			logger.ErrorContext(ctx, "got error on get-balance", slog.Any("error", err))
			return
		}
		if bal == nil {
			logger.ErrorContext(ctx, "get-balance return nil")
			return
		}

		// deposit money when needed
		if bal.Amount < amount {
			depositAmount := 100 + rand.Float64()*10000
			err := Deposit(ctx, balanceId, depositAmount)
			if err != nil {
				logger.ErrorContext(ctx, "got error on deposit", slog.Any("error", err))
				return
			}
		}

		// exec withdrawal
		err = Withdraw(ctx, balanceId, amount)
		if err != nil {
			logger.ErrorContext(ctx, "got error on withdrawal", slog.Any("error", err))
			return
		}
	})
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
