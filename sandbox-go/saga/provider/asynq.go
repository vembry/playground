package provider

import (
	"context"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type asynqSagaProvider struct {
	redisClientOpt asynq.RedisClientOpt
	client         *asynq.Client
	server         *asynq.Server
	servermux      *asynq.ServeMux
	queue          map[string]int
}

const (
	redisUrl = "redis://:local@host.docker.internal:6379/0"
)

func NewAsynqProvider() *asynqSagaProvider {
	// Single node Redis
	opt, _ := asynq.ParseRedisURI(redisUrl)
	redisClientOpt := opt.(asynq.RedisClientOpt)

	return &asynqSagaProvider{
		redisClientOpt: redisClientOpt,
		client:         asynq.NewClient(redisClientOpt),
		servermux:      asynq.NewServeMux(),
		queue:          map[string]int{},
	}
}

func (asp *asynqSagaProvider) Start() {
	asp.server = asynq.NewServer(
		asp.redisClientOpt,
		asynq.Config{
			Queues: asp.queue,
		},
	)

	go func() {
		if err := asp.server.Run(asp.servermux); err != nil {
			log.Fatalf("error on starting asynq server on saga. err=%v", err)
		}
	}()
}

func (asp *asynqSagaProvider) Close() {
	asp.client.Close()
	if asp.server != nil {
		asp.server.Shutdown()
	}
}

func (asp *asynqSagaProvider) Register(name string, taskCh chan string) {
	asp.servermux.HandleFunc(name, func(ctx context.Context, t *asynq.Task) error {
		taskCh <- time.Now().Format(time.RFC3339Nano)
		return nil
	})
	asp.queue[name] = 1
}

func (asp *asynqSagaProvider) GetAsynqClient() *asynq.Client {
	return asp.client
}
