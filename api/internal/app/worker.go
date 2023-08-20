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

// WithWorkers is to register individual workers into the worker
func (w *Worker) WithWorkers(workers ...cmd.WorkerHandler) {
	if len(workers) == 0 {
		log.Fatal("missing worker")
	}

	// setup worker's handlers
	for _, worker := range workers {
		worker := worker
		w.mux.HandleFunc(worker.Type(), worker.Perform)
	}
}

// WithMiddleware is to register middlewares to the worker
func (w *Worker) WithMiddleware(middlewares ...func(h asynq.Handler) asynq.Handler) {
	if len(middlewares) == 0 {
		return
	}

	// setup worker's middlewares
	for _, middleware := range middlewares {
		middleware := middleware
		w.mux.Use(middleware)
	}
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
