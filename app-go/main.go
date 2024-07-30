package main

import (
	"app/cmd"
	"app/internal/app"
	"app/internal/domain"
	internalhttp "app/internal/http"
	"app/internal/repository/postgres"
	repoRedis "app/internal/repository/redis"
	workerasynq "app/internal/worker/asynq"
	workerrabbit "app/internal/worker/rabbit"
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

	// setup worker
	workerAsynq := workerasynq.New(appConfig.RedisUri) // TODO: remove
	workerRabbit := workerrabbit.New(appConfig.RabbitUri)

	// setup individual asynq workers
	withdrawalWorkerAsynq := workerasynq.NewWithdrawal(workerAsynq.GetClient()) // TODO: remove
	depositWorkerAsynq := workerasynq.NewDeposit(workerAsynq.GetClient())       // TODO: remove
	transferWorkerAsynq := workerasynq.NewTransfer(workerAsynq.GetClient())     // TODO: remove

	// TODO: remove
	// register individual-workers to the asynq
	workerAsynq.RegisterWorkers(
		withdrawalWorkerAsynq,
		depositWorkerAsynq,
		transferWorkerAsynq,
	)

	// setup individual rabbit workers
	transferWorkerRabbit := workerrabbit.NewTransfer(workerRabbit.GetConnection())
	depositWorkerRabbit := workerrabbit.NewDeposit(workerRabbit.GetConnection())
	withdrawWorkerRabbit := workerrabbit.NewWithdraw(workerRabbit.GetConnection())

	// register individual-workers to the rabbit
	workerRabbit.RegisterWorkers(
		withdrawWorkerRabbit,
		depositWorkerRabbit,
		transferWorkerRabbit,
	)

	// setup domain(s)
	balanceDomain := domain.NewBalance(
		balanceRepo,
		depositRepo,
		withdrawalRepo,
		transferRepo,
		depositWorkerRabbit,
		withdrawWorkerRabbit,
		transferWorkerRabbit,
		ledgerRepo,
		lockerRepo,
	)

	// inject missing deps
	withdrawalWorkerAsynq.InjectDeps(balanceDomain) // TODO: remove
	depositWorkerAsynq.InjectDeps(balanceDomain)    // TODO: remove
	transferWorkerAsynq.InjectDeps(balanceDomain)   // TODO: remove

	withdrawWorkerRabbit.InjectDeps(balanceDomain)
	depositWorkerRabbit.InjectDeps(balanceDomain)
	transferWorkerRabbit.InjectDeps(balanceDomain)

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
		cmd.NewWork(
			appMetric,
			workerRabbit,
		),
		cmd.NewDummy(),
	)

	if err := cli.Execute(); err != nil {
		log.Fatalf("found error on executing app's cli. err=%v", err)
	}
}
