package asynq

import (
	"context"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

type asynqx struct {
	server          *asynq.Server
	client          *asynq.Client
	mux             *asynq.ServeMux
	redisOpt        asynq.RedisClientOpt
	queuePriorities map[string]int
}

func logging(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		log.Printf("processing task=%s. payload=%s", t.Type(), string(t.Payload()))
		err := h.ProcessTask(ctx, t)
		log.Printf("processed task=%s. payload=%s", t.Type(), string(t.Payload()))

		return err
	})
}

func New(redisUri string) *asynqx {
	opt, _ := asynq.ParseRedisURI(redisUri)
	redisOpt := opt.(asynq.RedisClientOpt)

	mux := asynq.NewServeMux()

	// middleware
	mux.Use(logging)

	return &asynqx{
		mux:             mux,
		redisOpt:        redisOpt,
		queuePriorities: make(map[string]int),
		client:          asynq.NewClient(redisOpt),
	}
}

func (w *asynqx) GetClient() *asynq.Client {
	return w.client
}

type IConsumer interface {
	Path() string
	Priority() int
	Handle(ctx context.Context, task *asynq.Task) error
}

func (w *asynqx) Name() string {
	return "asynq"
}

func (w *asynqx) RegisterWorker(consumers ...IConsumer) {
	for i := range consumers {
		// assign queue-priority
		w.queuePriorities[consumers[i].Path()] = consumers[i].Priority()

		// assign mux handler
		w.mux.HandleFunc(
			consumers[i].Path(),
			consumers[i].Handle,
		)
	}
}

func (w *asynqx) Start() {
	w.server = asynq.NewServer(
		w.redisOpt,
		asynq.Config{
			Concurrency:     50,
			Queues:          w.queuePriorities,
			ShutdownTimeout: 20 * time.Second,
			StrictPriority:  true,
		},
	)

	// starts worker
	if err := w.server.Start(w.mux); err != nil {
		log.Fatalf("error on starting asynq server. err=%v", err)
	}
}

func (w *asynqx) Stop() {
	if w.server != nil {
		w.server.Shutdown()
	}
	if w.client != nil {
		w.client.Close()
	}
}
