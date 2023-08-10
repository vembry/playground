package app

import (
	"api/cmd"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

// Worker contain the config of app-worker
type Worker struct {
	redisConn         *asynq.RedisClientOpt
	mux               *asynq.ServeMux
	server            *asynq.Server
	client            *asynq.Client
	queues            map[string]int
	postStartCallback func()
}

// NewServer is to setup app-server
func NewWorker(cfg *EnvConfig) *Worker {
	redisConnOpt, err := asynq.ParseRedisURI(cfg.RedisUri)
	if err != nil {
		log.Fatalf("failed to parse redis-uri. err=%v", err)
	}

	r, ok := redisConnOpt.(asynq.RedisClientOpt)
	if !ok {
		log.Fatalf("failed to convert redis-opt to redis-client-opt. err=%v", err)
	}

	mux := asynq.NewServeMux()
	mux.Use(func(h asynq.Handler) asynq.Handler {
		return logging(h)
	})

	return &Worker{
		redisConn: &r,
		mux:       asynq.NewServeMux(),
		queues:    make(map[string]int),
	}
}

// WithPostStartCallback inject callback to the post start callback
func (w *Worker) WithPostStartCallback(callback func()) {
	w.postStartCallback = callback
}

// Logging is the middleware that log-out the asynchrounous background task.
func logging(h asynq.Handler) asynq.Handler {
	return asynq.HandlerFunc(func(ctx context.Context, t *asynq.Task) error {
		log.Printf("processing '%s' background task with '%s'...", t.Type(), string(t.Payload()))
		err := h.ProcessTask(ctx, t)
		log.Printf("processed '%s' background task...", t.Type())
		return err
	})
}

// Start is to start worker
func (w *Worker) Start() error {
	// establish the worker
	w.server = asynq.NewServer(w.redisConn, asynq.Config{
		Concurrency:     10,
		Queues:          w.queues,
		ShutdownTimeout: 10 * time.Second,
	})

	// establish connection to the queue
	w.client = asynq.NewClient(w.redisConn)

	return w.server.Start(w.mux)
}

// Shutdown is to shutdown worker gracefully
func (w *Worker) Shutdown() error {
	// close connection to the queue
	w.client.Close()

	// signals worker to stop picking up queues
	w.server.Stop()
	w.server.Shutdown()

	return nil
}

// RegisterWorkers is to register individual workers into the app-worker
func (w *Worker) RegisterWorkers(workers ...cmd.WorkerHandler) error {
	if w.server == nil {
		return fmt.Errorf("missing worker's initialization")
	}

	if len(workers) == 0 {
		return fmt.Errorf("missing worker")
	}

	for _, worker := range workers {
		_worker := worker
		w.mux.HandleFunc(worker.Type(), func(ctx context.Context, task *asynq.Task) error {
			return _worker.Perform(ctx, task)
		})
	}

	return nil
}

// Enqueue is to enqueue task into the worker
func (w *Worker) Enqueue(ctx context.Context, task *asynq.Task, taskOptions ...asynq.Option) (*asynq.TaskInfo, error) {
	taksInfo, err := w.client.EnqueueContext(ctx, task, taskOptions...)
	if err != nil {
		return nil, fmt.Errorf("found error on enqueuing task to worker. err=%w", err)
	}

	return taksInfo, nil
}

// RegisterQueues is to register individual queues and it's respective priorities
func (w *Worker) RegisterQueues(queues map[string]int) {
	w.queues = queues
}

// ConnectToQueue establishes the connection to the Connection queue.
func (w *Worker) ConnectToQueue() {
	w.client = asynq.NewClient(w.redisConn)
}

// DisconnectFromQueue closes the connection to the Connection queue.
func (w *Worker) DisconnectFromQueue() error {
	if w.client != nil {
		return w.client.Close()
	}

	return nil
}

// GetPostStartCallback is to return post-start's callback
func (w *Worker) GetPostStartCallback() func() {
	return w.postStartCallback
}
