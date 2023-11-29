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
	config := app.NewConfig(embedFS)

	db, closer := app.NewOrmDb(config)
	defer closer() // close connection when main.go closes

	// setup repository(s)
	balanceRepository := postgres.NewBalance(db)
	// ledgerRepository := postgres.NewLedger(db)
	depositRepository := postgres.NewDeposit(db)
	withdrawalRepository := postgres.NewWithdrawal(db)
	transferRepository := postgres.NewTransfer(db)

	// setup domain(s)
	balanceDomain := domain.NewBalance(
		balanceRepository,
		depositRepository,
		withdrawalRepository,
		transferRepository,
	)

	withdrawalWorker := workerkafka.NewWithdrawal()

	// initiate CLI(s)
	cli := &cobra.Command{}
	cli.AddCommand(
		cmdserve.New(config, balanceDomain),
		cmdwork.New(withdrawalWorker),
	)

	if err := cli.Execute(); err != nil {
		log.Fatalf("found error on executing app's cli. err=%v", err)
	}
}
