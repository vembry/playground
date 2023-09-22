package main

import (
	"embed"
	"log"
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

	db, err := appDb.DB()
	if err != nil {
		log.Fatalf("found error on getting db instance. err=%v", err)
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.SetMaxOpenConns(100)

	// setup domain
	mutex := mutexDomain.New(appCache.GetClient())
	balance := balanceDomain.New(appDb, mutex)
	transaction := transactionDomain.New(appDb)

	// initiate individual worker
	pendingTransactionWorker := worker.NewPendingTransaction(transaction)
	addBalanceWorker := worker.NewAddBalance(balance)

	// setup app-worker
	appWorker := app.NewWorker(appConfig, map[string]int{
		pendingTransactionWorker.Queue(): 3,
		addBalanceWorker.Queue():         7,
	})
	appWorker.WithMiddleware(
		appMetric.AsynqTask,
	)
	appWorker.WithWorkers(
		pendingTransactionWorker,
		addBalanceWorker,
	)
	appWorker.WithPostStartCallback(func() {
		appMetric.StartServer()
	})
	appWorker.WithPostShutdownCallback(func() {
		appMetric.ShutdownServer()
	})

	// plug missing dependecies to transaction domain
	transaction.WithBalance(balance)
	transaction.WithPendingTransactionHandler(pendingTransactionWorker)

	// plug missing worker to worker-handlers
	pendingTransactionWorker.WithWorker(appWorker)
	addBalanceWorker.WithWorker(appWorker)

	// setup server's http-handler
	r := handler.NewHttpHandler(transaction, balance, addBalanceWorker, appMetric)

	// setup app-server
	appServer := app.NewServer(appConfig, r)
	appServer.WithPostStartCallback(func() {
		appWorker.ConnectToQueue()
		appMetric.StartServer()
	})
	appServer.WithPostShutdownCallback(func() {
		if err := appWorker.DisconnectFromQueue(); err != nil {
			log.Printf("found error on disconnecting from queue. err=%v", err)
		}
		appMetric.ShutdownServer()
	})

	// setup app-cli
	appCli := app.NewCli(appServer, appWorker)

	// start app cli
	if err := appCli.Execute(); err != nil {
		os.Exit(1)
	}
}
