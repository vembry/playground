package main

import (
	"embed"
	"os"

	"api/internal/app"
	"api/internal/app/handler"
	balanceDomain "api/internal/domain/balance"
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
	db, close := app.NewOrmDb(appConfig)
	// when main stack closes, then close db connection
	defer close()

	// setup domain
	balance := balanceDomain.New(db)
	transaction := transactionDomain.New(db)

	// initiate individual worker
	pendingTransactionWorker := worker.NewPendingTransaction(transaction)
	addBalanceWorker := worker.NewAddBalance(balance)

	// setup server's http-handler
	r := handler.NewHttpHandler(transaction, balance, addBalanceWorker)

	// setup app-server
	appServer := app.NewServer(appConfig, r)

	// setup app-worker
	appWorker := app.NewWorker(appConfig)
	appWorker.WithPostStartCallback(func() {
		// register individual workers to the app-worker
		appWorker.RegisterWorkers(
			pendingTransactionWorker,
			addBalanceWorker,
		)
	})

	// register individual queues + respective priority to the app-worker
	appWorker.RegisterQueues(map[string]int{
		pendingTransactionWorker.Queue(): 1,
		addBalanceWorker.Queue():         2,
	})

	// plug missing dependecies to transaction domain
	transaction.WithBalance(balance)
	transaction.WithPendingTransactionHandler(pendingTransactionWorker)

	// plug missing worker to worker-handlers
	pendingTransactionWorker.WithWorker(appWorker)
	addBalanceWorker.WithWorker(appWorker)

	// setup app-cli
	appCli := app.NewCli(appServer, appWorker)

	// start app cli
	if err := appCli.Execute(); err != nil {
		os.Exit(1)
	}
}
