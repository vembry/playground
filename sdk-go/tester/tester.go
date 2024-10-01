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

	Type     LoadType      // determine type of load, time based or counter based
	Duration time.Duration // max duration of load test for time based
	Counter  int           // max count of load test for counter based
}

// New initiate load-tester instance
func New(cfg Config) *tester {
	setDefault(&cfg) // handle defaults

	return &tester{
		cfg: cfg,
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
	if cfg.Counter == 0 {
		cfg.Counter = 1000
	}

	// when 'ConcurrentWorkerCount' not defined,
	// then defaults to 1
	if cfg.ConcurrentWorkerCount == 0 {
		cfg.ConcurrentWorkerCount = 1
	}
}

// Do setup and execute load-tester
func (t *tester) Do(execution func(ctx context.Context, l *slog.Logger)) {
	// load otel trace
	tracer := otel.Tracer("load-test-tracer")

	// setup array of channels as task executor
	taskChannels := make([]chan int, t.cfg.ConcurrentWorkerCount)
	for i := range taskChannels {
		taskChannels[i] = make(chan int, 100)
	}
	defer func() {
		// close channels when it's done
		for i := range taskChannels {
			close(taskChannels[i])
		}
	}()

	taskExecutorWg := sync.WaitGroup{} // this is to control the task executor
	taskProducerWg := sync.WaitGroup{} // this is to control the task producer

	t.cfg.Logger.Enabled(context.Background(), slog.LevelInfo)

	// setup concurrent worker to execute task
	for i, taskChannel := range taskChannels {
		// run scenario in goroutine
		go func(channelId int, taskCh chan int) {
			t.cfg.Logger.Debug("starting task worker", slog.Int("channelId", channelId))

			// load task
			for taskId := range taskCh {
				t.cfg.Logger.Debug("processing task", slog.Int("channelId", channelId), slog.Int("taskId", taskId))
				t.do(tracer, execution) // execute scenario
				taskExecutorWg.Done()
			}

			t.cfg.Logger.Debug("stopping task worker", slog.Int("channelId", channelId))
		}(i, taskChannel)
	}

	// prep load conditions
	start := time.Now()
	end := start.Add(t.cfg.Duration)
	counter := 0

	t.cfg.Logger.Info("starting load test...")

	// start actual load test
	for t.isConditionAllow(counter, end) {
		taskExecutorWg.Add(len(taskChannels)) // add count to wait-group

		for _, taskCh := range taskChannels {
			taskProducerWg.Add(1)

			counter += 1
			taskCh <- counter // enqueue task using goroutine

			if counter%1000 == 0 {
				// log every '1000' count
				t.cfg.Logger.Debug("counter checkpoint", slog.Int("counter", counter))
			}
			taskProducerWg.Done()
		}
	}

	// waiting for concurrent worker to finish
	t.cfg.Logger.Info("looper finished...", slog.String("duration", time.Since(start).String()), slog.Int("counter", counter))
	taskProducerWg.Wait()
	taskExecutorWg.Wait()

	t.cfg.Logger.Info("worker finished", slog.String("duration", time.Since(start).String()), slog.Int("counter", counter)) // NOTE:  peak at 31,042,670,712 counts per 5 mins without operation
}

// do executes test scenario
func (t *tester) do(tracer trace.Tracer, scenario func(ctx context.Context, l *slog.Logger)) {
	// start open-telemetry
	ctx, span := tracer.Start(
		context.Background(),
		"load-test-starter",
	)
	defer span.End()

	// execute test script
	scenario(ctx, t.cfg.Logger)
}

// isConditionAllow validates whether task-producer are still allowed to produce
func (t *tester) isConditionAllow(counter int, end time.Time) bool {
	switch t.cfg.Type {
	case LoadType_Count:
		return counter < t.cfg.Counter
	case LoadType_Time:
		return time.Now().Before(end)
	}
	return false
}
