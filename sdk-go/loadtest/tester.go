package loadtest

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type tester struct {
	cfg Config

	tracer trace.Tracer
}

// LoadType specify load test type
type LoadType int

const (
	LoadType_Time LoadType = iota + 1
	LoadType_Count
)

// Config contain load-tester configuration
type Config struct {
	Logger      *slog.Logger
	WorkerCount int

	Type       LoadType      // determine type of load, time based or counter based
	Duration   time.Duration // max duration of load test for time based
	MaxCounter int           // max count of load test for counter based
}

// New initiate load-tester instance
func New(cfg Config) *tester {
	setDefault(&cfg) // handle defaults

	return &tester{
		cfg:    cfg,
		tracer: otel.Tracer("load-test"),
	}
}

// Do setup and execute load-tester
func (t *tester) Do(scenario func(ctx context.Context, l *slog.Logger)) {
	// setup task channel
	taskChannel := make(chan int, 100) // setting channel buffer to 100
	defer func() {
		close(taskChannel)
	}()

	taskExecutorWg := sync.WaitGroup{}

	// run task's executor
	go t.taskExecutor(taskChannel, func() {
		t.do(scenario)
		taskExecutorWg.Done()
	})

	// load test execution pre-reqs
	start := time.Now()
	counter := 0

	// produce tasks
	for t.isConditionAllow(counter, start.Add(t.cfg.Duration)) {
		taskExecutorWg.Add(1)
		counter++
		taskChannel <- counter
	}

	// wait for task's executor to finish
	t.cfg.Logger.Info("awaiting worker...")
	taskExecutorWg.Wait()

	// post load test infos
	t.cfg.Logger.Info(
		"worker finished",
		slog.Int("counter", counter),
		slog.String("duration", time.Since(start).String()),
	)

	// publish metric here?
	// ...
}

// do executes test scenario
func (t *tester) do(scenario func(ctx context.Context, l *slog.Logger)) {
	// start open-telemetry
	ctx, span := t.tracer.Start(
		context.Background(),
		"load-test-start",
	)

	// close span on exiting function
	defer span.End()

	// execute test script
	scenario(ctx, t.cfg.Logger)

	// publish metric here?
	// ...
}

// isConditionAllow validates whether task-producer are still allowed to produce
func (t *tester) isConditionAllow(counter int, timeLimit time.Time) bool {
	switch t.cfg.Type {
	case LoadType_Count:
		return counter < t.cfg.MaxCounter
	case LoadType_Time:
		return time.Now().Before(timeLimit)
	}
	return false
}

// taskExecutor encapsulate concurrent executor
func (t *tester) taskExecutor(taskChannel chan int, callback func()) {
	for i := range t.cfg.WorkerCount {
		go func(workerId int) {
			t.cfg.Logger.Debug("starting worker", slog.Int("worker", workerId))
			for task := range taskChannel {
				t.cfg.Logger.Info("processing task", slog.Int("worker", workerId), slog.Int("task", task))

				callback()
			}
			t.cfg.Logger.Debug("shutting down worker", slog.Int("worker", workerId))
		}(i)
	}
}
