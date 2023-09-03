package main

import (
	"embed"
	"os"

	"app-go/internal/app"
	balanceDomain "app-go/internal/domain/balance"
	mutexDomain "app-go/internal/domain/mutex"
	transactionDomain "app-go/internal/domain/transaction"
	"app-go/internal/handler"
	"app-go/internal/worker"
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

	// setup app-metric
	appMetric := app.NewMetric(appConfig)

	// setup db
	appDb, closeCallback := app.NewOrmDb(appConfig)
	// when main stack closes, then close db connection
	defer closeCallback()

	// setup domain
	mutex := mutexDomain.New(appCache.GetClient())
	balance := balanceDomain.New(appDb, mutex)
	transaction := transactionDomain.New(appDb)

	// initiate individual worker
	pendingTransactionWorker := worker.NewPendingTransaction(transaction)
	addBalanceWorker := worker.NewAddBalance(balance)

	// setup server's http-handler
	r := handler.NewHttpHandler(transaction, balance, addBalanceWorker, appMetric)

	// setup app-server
	appServer := app.NewServer(appConfig, r)
	appServer.WithPostStartCallback(func() {
		// start metric server
		appMetric.Start()
	})

	// setup app-worker
	appWorker := app.NewWorker(appConfig)
	appWorker.WithMiddleware(
		appMetric.AsynqTask,
	)
	appWorker.WithWorkers(
		pendingTransactionWorker,
		addBalanceWorker,
	)
	appWorker.WithPostStartCallback(func() {
		// start metric server
		appMetric.Start()

		// // register individual workers to the app-worker
		// appWorker.RegisterWorkers(
		// 	pendingTransactionWorker,
		// 	addBalanceWorker,
		// )
	})

	// register individual queues + respective priority to the app-worker
	appWorker.RegisterQueues(map[string]int{
		pendingTransactionWorker.Queue(): 3,
		addBalanceWorker.Queue():         7,
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
