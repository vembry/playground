package main

import (
	"app/cmd"
	"app/internal/app"
	balanceModuleRepo "app/internal/module/balance/repository/postgres"
	balancehttpserver "app/internal/module/balance/server/http"
	balanceModuleService "app/internal/module/balance/service"
	"app/internal/module/locker"
	"app/internal/server/http"
	"app/internal/telemetry"
	"app/internal/worker/dummy"
	workerrabbit "app/internal/worker/rabbit"
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
	// setup app's pre-requisites
	// ==========================

	// setup config
	appConfig := app.NewConfig(embedFS)

	// setup telemetry
	telemetryShutdown := telemetry.New()
	defer telemetryShutdown()

	// setup cache
	cacheOpts, err := redis.ParseURL(appConfig.RedisUri)
	if err != nil {
		log.Fatalf("failed to parse redis url. err=%v", err)
	}
	appCache := redis.NewClient(cacheOpts)
	defer appCache.Close()

	// setup db
	appDb, closer := app.NewOrmDb(appConfig)
	defer closer() // close connection when main.go closes

	// setup internal modules
	// ======================

	// setup repository
	balanceRepo := balanceModuleRepo.NewBalance(appDb)
	ledgerRepo := balanceModuleRepo.NewLedger(appDb)
	depositRepo := balanceModuleRepo.NewDeposit(appDb)
	withdrawalRepo := balanceModuleRepo.NewWithdrawal(appDb)
	transferRepo := balanceModuleRepo.NewTransfer(appDb)

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
	lockermodule := locker.New(appCache)
	balancemodule := balanceModuleService.New(
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
	balancehttp := balancehttpserver.New(balancemodule)

	// setup http server
	httpserver := http.New(appConfig.HttpAddress, balancehttp.GetHandler())

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
