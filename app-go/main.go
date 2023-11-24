package main

import (
	servecmd "app/cmd/serve"
	workcmd "app/cmd/work"
	"app/internal/app"
	"app/internal/domain"
	"app/internal/repository/postgres"
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

	appDb, closer := app.NewOrmDb(appConfig)
	defer closer() // close connection when main.go closes

	// setup repository(s)
	balanceRepository := postgres.NewBalance(appDb)
	ledgerRepository := postgres.NewLedger(appDb)

	// setup domain(s)
	balanceDomain := domain.NewBalance(balanceRepository, ledgerRepository)

	// initiate CLI(s)
	appCli := &cobra.Command{}
	appCli.AddCommand(
		servecmd.New(appConfig, balanceDomain),
		workcmd.New(),
	)

	if err := appCli.Execute(); err != nil {
		log.Fatalf("found error on executing app's cli. err=%v", err)
	}
}
