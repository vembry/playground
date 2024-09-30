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
	// handle defaults

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

	return &tester{
		cfg: cfg,
	}
}

// Do setup and execute load-tester
func (t *tester) Do(execution func(ctx context.Context, l *slog.Logger)) {

	// load otel trace
	tracer := otel.Tracer("load-test-tracer")

	// setting up array of channel
	taskChannels := make([]chan int, t.cfg.ConcurrentWorkerCount) // TODO: this should be parameterized
	for i := range taskChannels {
		taskChannels[i] = make(chan int)
	}
	defer func() {
		for i := range taskChannels {
			close(taskChannels[i])
		}
	}()

	taskWg := sync.WaitGroup{}

	t.cfg.Logger.Enabled(context.Background(), slog.LevelInfo)

	// setup concurrent worker to run the scenario
	for i, taskChannel := range taskChannels {
		// run scenario in goroutine
		go func(channelId int, taskCh chan int) {
			t.cfg.Logger.Debug("starting task worker", slog.Int("channelId", channelId))

			// load data
			for taskId := range taskCh {
				t.cfg.Logger.Debug("processing task", slog.Int("channelId", channelId), slog.Int("taskId", taskId))
				t.do(tracer, execution)
				taskWg.Done()
			}
			t.cfg.Logger.Debug("stopping task worker", slog.Int("channelId", channelId))
		}(i, taskChannel)
	}

	timeCounterPivot := time.Now() // for logging purposes

	t.cfg.Logger.Info("starting")

	looperWg := sync.WaitGroup{}

	// prep load conditions
	start := time.Now()
	end := start.Add(t.cfg.Duration)
	counter := 0

	// start actual load test
	for t.isConditionAllow(counter, end) {
		taskWg.Add(len(taskChannels)) // add count to wait-group

		for _, taskCh := range taskChannels {
			looperWg.Add(1)

			// enqueue task using goroutine because enqueuing
			// task to a channel directly has some latency. idk
			go func() {
				counter++
				if counter%1000 == 0 {
					// log every '1000' count
					t.cfg.Logger.Info("counter checkpoint", slog.Int("counter", counter), slog.String("duration", time.Since(timeCounterPivot).String()))
					timeCounterPivot = time.Now()
				}
				taskCh <- counter
				looperWg.Done()
			}()
		}
	}

	// waiting for concurrent worker to finish
	t.cfg.Logger.Info("looper finished...", slog.String("duration", time.Since(start).String()))
	looperWg.Wait()
	taskWg.Wait()

	t.cfg.Logger.Info("finished", slog.Int("counter", counter)) // NOTE:  peak at 31,042,670,712 counts per 5 mins without operation
}

// do executes tests
func (t *tester) do(tracer trace.Tracer, execution func(ctx context.Context, l *slog.Logger)) {
	// start open-telemetry
	ctx, span := tracer.Start(
		context.Background(),
		"load-test-starter",
	)
	defer span.End()

	// execute test script
	execution(ctx, t.cfg.Logger)
}

func (t *tester) isConditionAllow(counter int, end time.Time) bool {
	switch t.cfg.Type {
	case LoadType_Count:
		return counter < t.cfg.Counter
	case LoadType_Time:
		return time.Now().After(end)
	}

	return false
}
