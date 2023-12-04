package main

import (
	"app/cmd"
	"app/internal/app"
	"app/internal/domain"
	internalhttp "app/internal/http"
	"app/internal/repository/postgres"
	workerasynq "app/internal/worker/asynq"
	"embed"
	"log"

	"github.com/spf13/cobra"
)

var (
	//go:embed configs
	embedFS embed.FS
)

func main() {
	// setup config
	appConfig := app.NewConfig(embedFS)

	// setup app metric
	appMetric := app.NewMetric(appConfig)

	// setup db
	appDb, closer := app.NewOrmDb(appConfig)
	defer closer() // close connection when main.go closes

	// setup repository(s)
	balanceRepository := postgres.NewBalance(appDb)
	// ledgerRepository := postgres.NewLedger(appDb)
	depositRepository := postgres.NewDeposit(appDb)
	withdrawalRepository := postgres.NewWithdrawal(appDb)
	transferRepository := postgres.NewTransfer(appDb)

	// setup asynq worker
	workerAsynq := workerasynq.New(appConfig.RedisUri)

	// setup individual asynq workers
	withdrawalWorker := workerasynq.NewWithdrawal(workerAsynq.GetClient())

	// register individual-workers to the asynq
	workerAsynq.RegisterWorker(withdrawalWorker)

	// setup domain(s)
	balanceDomain := domain.NewBalance(
		balanceRepository,
		depositRepository,
		withdrawalRepository,
		transferRepository,
		withdrawalWorker,
	)

	// inject missing deps
	withdrawalWorker.InjectDeps(balanceDomain)

	httpserver := internalhttp.NewServer(
		appConfig.HttpAddress,
		appMetric,
		balanceDomain,
	)

	// initiate CLI(s)
	cli := &cobra.Command{}
	cli.AddCommand(
		cmd.NewServe(
			httpserver,
			appMetric,
		),
		cmd.NewWork(workerAsynq),
	)

	if err := cli.Execute(); err != nil {
		log.Fatalf("found error on executing app's cli. err=%v", err)
	}
}
