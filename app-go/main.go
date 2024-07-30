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

	tracerShutdownHandler := app.NewTracer()
	defer tracerShutdownHandler()

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

	// setup asynq
	workerAsynq := workerasynq.New(appConfig.RedisUri)

	// setup individual asynq workers
	withdrawalWorkerAsynq := workerasynq.NewWithdrawal(workerAsynq.GetClient())
	depositWorkerAsynq := workerasynq.NewDeposit(workerAsynq.GetClient())
	transferWorkerAsynq := workerasynq.NewTransfer(workerAsynq.GetClient())

	// register individual-workers to the asynq
	workerAsynq.RegisterWorkers(
		withdrawalWorkerAsynq,
		depositWorkerAsynq,
		transferWorkerAsynq,
	)

	// setup domain(s)
	balanceDomain := domain.NewBalance(
		balanceRepo,
		depositRepo,
		withdrawalRepo,
		transferRepo,
		depositWorkerAsynq,
		withdrawalWorkerAsynq,
		transferWorkerAsynq,
		ledgerRepo,
		lockerRepo,
	)

	// inject missing deps
	withdrawalWorkerAsynq.InjectDeps(balanceDomain)
	depositWorkerAsynq.InjectDeps(balanceDomain)
	transferWorkerAsynq.InjectDeps(balanceDomain)

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
		cmd.NewWork(appMetric, workerAsynq),
		cmd.NewDummy(),
	)

	if err := cli.Execute(); err != nil {
		log.Fatalf("found error on executing app's cli. err=%v", err)
	}
}
