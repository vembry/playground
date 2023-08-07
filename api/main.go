package main

import (
	"embed"
	"log"
	"os"

	"api/internal/app"
	"api/internal/app/handler"
	ledgerDomain "api/internal/domain/ledger"
	transactionDomain "api/internal/domain/transaction"
	"api/internal/worker"
)

var (
	//go:embed configs
	embedFS embed.FS
)

func main() {
	// setup config
	appConfig := app.NewConfig(embedFS)

	// setup db
	db := app.NewOrmDb(appConfig)

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

	// setup server's http-handler
	r := handler.NewHttpHandler(transaction, ledger)

	// setup app-server
	appServer := app.NewServer(appConfig, r.Handler())

	// initiate individual worker
	pendingTransactionWorker := worker.NewPendingTransaction(transaction)

	// setup app-worker
	appWorker := app.NewWorker(appConfig)
	appWorker.WithPostStartCallback(func() {
		// register individual workers to the app-worker
		appWorker.RegisterWorkers(
			pendingTransactionWorker,
		)
	})

	// register individual queues + respective priority to the app-worker
	appWorker.RegisterQueues(map[string]int{
		pendingTransactionWorker.Queue(): 1,
	})

	// plug missing dependecies to transaction domain
	transaction.WithLedger(ledger)
	transaction.WithPendingTransactionHandler(pendingTransactionWorker)

	// plug missing worker to worker-handler
	pendingTransactionWorker.WithWorker(appWorker)

	// setup app-cli
	appCli := app.NewCli(appServer, appWorker)

	// start app cli
	if err := appCli.Execute(); err != nil {
		os.Exit(1)
	}
}
