package main

import (
	"embed"
	"os"

	"api/internal/app"
	balanceDomain "api/internal/domain/balance"
	mutexDomain "api/internal/domain/mutex"
	transactionDomain "api/internal/domain/transaction"
	"api/internal/handler"
	"api/internal/worker"
)

var (
	//go:embed configs
	embedFS embed.FS
)

func main() {
	// setup app-config
	appConfig := app.NewConfig(embedFS)

	// setup app-cache
	appCache := app.NewCache(appConfig)

	// setup app-prometheus
	appPrometheus := app.NewPrometheus(appConfig)

	// setup db
	db, close := app.NewOrmDb(appConfig)
	// when main stack closes, then close db connection
	defer close()

	// setup domain
	mutex := mutexDomain.New(appCache.GetClient())
	balance := balanceDomain.New(db, mutex)
	transaction := transactionDomain.New(db)

	// initiate individual worker
	pendingTransactionWorker := worker.NewPendingTransaction(transaction)
	addBalanceWorker := worker.NewAddBalance(balance)

	// setup server's http-handler
	r := handler.NewHttpHandler(transaction, balance, addBalanceWorker)

	// setup app-server
	appServer := app.NewServer(appConfig, r)
	appServer.WithPostStartCallback(func() {
		// start prometheus server
		appPrometheus.Start()
	})

	// setup app-worker
	appWorker := app.NewWorker(appConfig)
	appWorker.WithPostStartCallback(func() {
		// start prometheus server
		appPrometheus.Start()

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
