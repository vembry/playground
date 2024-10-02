package tester

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
	Logger                *slog.Logger
	ConcurrentWorkerCount int

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

func setDefault(cfg *Config) {
	// default type to count when undefined
	if cfg.Type == 0 {
		cfg.Type = LoadType_Time
	}

	// default duration
	if cfg.Duration == 0 {
		cfg.Duration = 5 * time.Second
	}

	// default counter
	if cfg.MaxCounter == 0 {
		cfg.MaxCounter = 1000
	}

	// when 'ConcurrentWorkerCount' not defined,
	// then defaults to 1
	if cfg.ConcurrentWorkerCount == 0 {
		cfg.ConcurrentWorkerCount = 1
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
	counterValidator := 0 // to verify the of task executed vs task produced

	// setup task executor
	for workerId := range t.cfg.ConcurrentWorkerCount {
		go func(_workerId int) {
			t.cfg.Logger.Debug("starting worker", slog.Int("workerId", _workerId))
			for task := range taskChannel {
				t.cfg.Logger.Info("processing task", slog.Int("workerId", _workerId), slog.Int("task", task))

				t.do(scenario) // execute scenario

				taskExecutorWg.Done()
				counterValidator++
			}
			t.cfg.Logger.Debug("shutting down worker", slog.Int("workerId", _workerId))
		}(workerId)
	}

	// load test execution pre-reqs
	start := time.Now()
	counter := 0

	// run task producer
	for t.isConditionAllow(counter, start.Add(t.cfg.Duration)) {
		counter++
		taskChannel <- counter
		taskExecutorWg.Add(1)
	}

	// wait for concurrent worker to finish
	t.cfg.Logger.Info("awaiting worker...")
	taskExecutorWg.Wait()

	// post load test infos
	t.cfg.Logger.Info(
		"worker finished",
		slog.Int("counter", counter),
		slog.Int("counter_validator", counterValidator),
		slog.String("duration", time.Since(start).String()),
	)
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
}

// isConditionAllow validates whether task-producer are still allowed to produce
func (t *tester) isConditionAllow(counter int, start time.Time) bool {
	switch t.cfg.Type {
	case LoadType_Count:
		return counter < t.cfg.MaxCounter
	case LoadType_Time:
		return time.Now().Before(start)
	}
	return false
}
