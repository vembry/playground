package main

import (
	"embed"
	"log"
	"os"

	"api/cmd"
	"api/internal/app"
	"api/internal/app/handler"
	ledgerDomain "api/internal/domain/ledger"
	transactionDomain "api/internal/domain/transaction"

	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	//go:embed configs
	embedFS embed.FS
)

func main() {
	// setup config
	appConfig := app.NewConfig(embedFS)

	// setup db
	db := newOrmDb(appConfig)

	// close db connection on app closure
	defer func() {
		sqlDb, err := db.DB()
		if err != nil {
			log.Fatalf("found error on getting DB. err=%v", err)
		}

		err = sqlDb.Close()
		if err != nil {
			log.Fatalf("found error on closing DB connection. err=%v", err)
		}
	}()

	// setup ledger domain
	ledger := ledgerDomain.New(db)

	// setup transaction domain
	transaction := transactionDomain.New(db)
	transaction.WithLedger(ledger) // plug ledger to transaction domain

	// setup server's http-handler
	r := handler.NewHttpHandler(transaction, ledger)

	// setup app-server
	appServer := app.NewServer(appConfig, r.Handler())

	// start app-server
	appCli := newCli(appServer)

	if err := appCli.Execute(); err != nil {
		os.Exit(1)
	}
}

// newOrmDb is to initialize DB in ORM form using gorm
func newOrmDb(cfg *app.EnvConfig) *gorm.DB {
	// dsn := "host=localhost user=local password=local dbname=credits port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(cfg.DBConn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to initiate db. err=%v", err)
	}

	return db
}

// newCli is to construct clis
func newCli(server *app.Server) *cobra.Command {
	command := &cobra.Command{}
	command.AddCommand(cmd.NewServe(server))

	return command
}
