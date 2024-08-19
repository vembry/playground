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

// Config contain load-tester configuration
type Config struct {
	Logger                *slog.Logger
	Duration              time.Duration
	ConcurrentWorkerCount int
}

// New initiate load-tester instance
func New(cfg Config) *tester {
	// when not defined, defaulted to 1
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
	taskChannelArray := make([]chan int, 5) // TODO: this should be parameterized
	for i := range taskChannelArray {
		taskChannelArray[i] = make(chan int)
	}
	defer func() {
		for i := range taskChannelArray {
			close(taskChannelArray[i])
		}
	}()

	taskWg := sync.WaitGroup{}

	t.cfg.Logger.Enabled(context.Background(), slog.LevelInfo)

	// setup concurrent worker to run the scenario
	for channelIdx, taskCh := range taskChannelArray {

		// setup concurrent worker for each channel
		// from channel array
		for i := range t.cfg.ConcurrentWorkerCount {

			// run scenario in goroutine
			go func(channelIdx int, workerId int, taskCh chan int) {
				t.cfg.Logger.Debug("starting task worker", slog.Int("channelIdx", channelIdx), slog.Int("workerId", workerId))

				// load data
				for taskId := range taskCh {
					t.cfg.Logger.Debug("processing task", slog.Int("channelIdx", channelIdx), slog.Int("workerId", workerId), slog.Int("taskId", taskId))
					t.do(tracer, execution)
					taskWg.Done()
				}
				t.cfg.Logger.Debug("stopping task worker", slog.Int("channelIdx", channelIdx), slog.Int("workerId", workerId))
			}(channelIdx, i, taskCh)
		}
	}

	// setup basic counters
	counter := 0
	counterCheckpoint := 1000
	timeCounterPivot := time.Now()
	start := time.Now()

	t.cfg.Logger.Info("starting")

	looperWg := sync.WaitGroup{}

	// start enqueuing task to the concurrent worker
	for counter < 60000 { // TODO: this should be parameterized
		taskWg.Add(len(taskChannelArray))

		for _, taskCh := range taskChannelArray {
			looperWg.Add(1)

			// enqueue task using go routine
			// because enqueuing task to a channel
			// has some latency
			go func() {
				counter++
				if counter%counterCheckpoint == 0 {
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
