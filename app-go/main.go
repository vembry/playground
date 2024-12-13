package main

import (
	"app/cmd"
	"app/internal/app"
	"app/internal/module/balance"
	"app/internal/module/locker"
	"app/internal/repository/postgres"
	"app/internal/server/http"
	"app/internal/server/http/handler"
	"app/internal/worker/dummy"
	workerrabbit "app/internal/worker/rabbit"
	"context"
	"embed"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cobra"
)

var (
	//go:embed configs
	embedFS embed.FS
)

func main() {
	ctx := context.Background()

	// setup app's pre-requisites
	// ==========================

	// setup config
	appConfig := app.NewConfig(embedFS)

	// setup telemetry
	telemetryShutdownHandler, err := app.NewTelemetry(ctx)
	if err != nil {
		log.Fatalf("failed to initiate telemetry")
	}
	defer telemetryShutdownHandler()

	// setup cache
	cacheOpts, err := redis.ParseURL(appConfig.RedisUri)
	if err != nil {
		log.Fatalf("failed to parse redis url. err=%v", err)
	}
	cache := redis.NewClient(cacheOpts)
	defer cache.Close()

	// setup db
	appDb, closer := app.NewOrmDb(appConfig)
	defer closer() // close connection when main.go closes

	// setup internal modules
	// ======================

	// setup repository
	balanceRepo := postgres.NewBalance(appDb)
	ledgerRepo := postgres.NewLedger(appDb)
	depositRepo := postgres.NewDeposit(appDb)
	withdrawalRepo := postgres.NewWithdrawal(appDb)
	transferRepo := postgres.NewTransfer(appDb)

	// setup worker
	workerRabbit := workerrabbit.New(appConfig.RabbitUri)
	workerDummy := dummy.New()

	// setup individual rabbit workers
	transferWorkerRabbit := workerrabbit.NewTransfer(workerRabbit.GetConnection())
	depositWorkerRabbit := workerrabbit.NewDeposit(workerRabbit.GetConnection())
	withdrawWorkerRabbit := workerrabbit.NewWithdraw(workerRabbit.GetConnection())

	// register individual-workers to the rabbit
	workerRabbit.RegisterWorkers(
		withdrawWorkerRabbit,
		depositWorkerRabbit,
		transferWorkerRabbit,
	)

	// setup modules
	lockermodule := locker.New()
	balancemodule := balance.New(
		balanceRepo,
		depositRepo,
		withdrawalRepo,
		transferRepo,
		depositWorkerRabbit,
		withdrawWorkerRabbit,
		transferWorkerRabbit,
		ledgerRepo,
		lockermodule,
	)

	// inject modules to workers
	withdrawWorkerRabbit.InjectDeps(balancemodule)
	depositWorkerRabbit.InjectDeps(balancemodule)
	transferWorkerRabbit.InjectDeps(balancemodule)

	// setup http server
	// =================

	// setup handler(s)
	balanceHandler := handler.NewBalance(balancemodule)

	// setup http server
	httpserver := http.New(appConfig.HttpAddress, balanceHandler.GetMux())

	// initiate CLI(s)
	cli := &cobra.Command{}
	cli.AddCommand(
		cmd.NewServe(
			httpserver,
		),
		cmd.NewWork(
			workerRabbit,
			workerDummy,
		),
		cmd.NewDummy(),
	)

	// run app
	if err := cli.Execute(); err != nil {
		log.Fatalf("found error on executing app's cli. err=%v", err)
	}
}
