package main

import (
	"app/cmd"
	"app/internal/app"
	"app/internal/domain"
	internalhttp "app/internal/http"
	"app/internal/repository/postgres"
	repoRedis "app/internal/repository/redis"
	workerasynq "app/internal/worker/asynq"
	"embed"
	"log"

	"github.com/redis/go-redis/v9"
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

	// setup cache
	cacheOpts, err := redis.ParseURL(appConfig.RedisUri)
	if err != nil {
		log.Fatalf("failed to parse redis url. err=%v", err)
	}
	cache := redis.NewClient(cacheOpts)
	defer cache.Close()

	// setup db
	appDb, closer := app.NewOrmDb(appConfig)
	defer closer() // close connection when main.go closes

	// setup repository(s)
	balanceRepo := postgres.NewBalance(appDb)
	ledgerRepo := postgres.NewLedger(appDb)
	depositRepo := postgres.NewDeposit(appDb)
	withdrawalRepo := postgres.NewWithdrawal(appDb)
	transferRepo := postgres.NewTransfer(appDb)
	lockerRepo := repoRedis.NewLocker(cache)

	// setup asynq worker
	workerAsynq := workerasynq.New(appConfig.RedisUri)

	// setup individual asynq workers
	withdrawalWorker := workerasynq.NewWithdrawal(workerAsynq.GetClient())
	depositWorker := workerasynq.NewDeposit(workerAsynq.GetClient())
	transferWorker := workerasynq.NewTransfer(workerAsynq.GetClient())

	// register individual-workers to the asynq
	workerAsynq.RegisterWorker(
		withdrawalWorker,
		depositWorker,
		transferWorker,
	)

	// setup domain(s)
	balanceDomain := domain.NewBalance(
		balanceRepo,
		depositRepo,
		withdrawalRepo,
		transferRepo,
		depositWorker,
		withdrawalWorker,
		transferWorker,
		ledgerRepo,
		lockerRepo,
	)

	// inject missing deps
	withdrawalWorker.InjectDeps(balanceDomain)
	depositWorker.InjectDeps(balanceDomain)
	transferWorker.InjectDeps(balanceDomain)

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
