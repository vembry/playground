package main

import (
	cmdserve "app/cmd/serve"
	cmdwork "app/cmd/work"
	"app/internal/app"
	"app/internal/domain"
	"app/internal/repository/postgres"
	workerkafka "app/internal/worker/kafka"
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

	// setup app-cache
	appCache := app.NewCache(appConfig)

	appDb, closer := app.NewOrmDb(appConfig)
	defer closer() // close connection when main.go closes

	// setup repository(s)
	balanceRepository := postgres.NewBalance(appDb)
	// ledgerRepository := postgres.NewLedger(appDb)
	depositRepository := postgres.NewDeposit(appDb)
	withdrawalRepository := postgres.NewWithdrawal(appDb)
	transferRepository := postgres.NewTransfer(appDb)

	// setup worker
	withdrawalWorker := workerkafka.NewWithdrawal(appConfig)

	// setup domain(s)
	mutexDomain := domain.NewMutex(appCache.GetClient())
	balanceDomain := domain.NewBalance(
		balanceRepository,
		depositRepository,
		withdrawalRepository,
		transferRepository,
		withdrawalWorker,
		mutexDomain,
	)

	// inject missing dependencies
	withdrawalWorker.InjectDep(balanceDomain)

	// initiate CLI(s)
	cli := &cobra.Command{}
	cli.AddCommand(
		cmdserve.New(appConfig, balanceDomain),
		cmdwork.New(withdrawalWorker),
	)

	if err := cli.Execute(); err != nil {
		log.Fatalf("found error on executing app's cli. err=%v", err)
	}
}
