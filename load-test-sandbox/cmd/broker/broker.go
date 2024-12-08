package broker

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"sdk/broker"
	sdkpb "sdk/broker/pb"
	"sync"

	loadtest "github.com/vembry/load-test"

	"github.com/segmentio/ksuid"
	"github.com/spf13/cobra"
)

func New(logger *slog.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "broker",
		Short: "load test execution for broker",
		Long:  "load test execution for broker long",
		Run: func(cmd *cobra.Command, args []string) {
			// setup load tester
			t := loadtest.New(
				loadtest.Config{
					Type:        loadtest.LoadType_Count,
					Logger:      logger,
					WorkerCount: 10,
					MaxCounter:  100000,
				},
			)

			// setup parameter for test
			params := []string{}
			for i := range 5 {
				params = append(params, fmt.Sprintf("queue_%d", i))
			}

			// setup broker client
			client, shutdown := broker.New("0.0.0.0:4000")
			defer shutdown()

			wg := sync.WaitGroup{}
			wg.Add(2)

			// run load tester

			// run enqueues
			go func() {
				t.Do(func(ctx context.Context, logger *slog.Logger) {
					// choose balance id
					i := randRange(0, len(params)-1)
					queueName := params[i]

					// enqueue
					_, err := client.GRPC().Enqueue(ctx, &sdkpb.EnqueueRequest{
						QueueName: queueName,
						Payload:   ksuid.New().String(),
					})
					if err != nil {
						logger.Error("'Enqueue' return error", slog.String("error", err.Error()))
					}
				})
				wg.Done()
			}()

			// run polls and ack
			go func() {
				t.Do(func(ctx context.Context, logger *slog.Logger) {
					// choose balance id
					i := randRange(0, len(params)-1)
					queueName := params[i]

					// poll it
					queue, err := client.GRPC().Poll(ctx, &sdkpb.PollRequest{
						QueueName: queueName,
					})
					if err != nil {
						logger.Error("'Poll' return error", slog.String("error", err.Error()))
					}

					// this means when theres no queue found
					if queue.Data == nil {
						return
					}

					// complete it
					_, err = client.GRPC().CompletePoll(ctx, &sdkpb.CompletePollRequest{
						QueueId: queue.Data.Id,
					})
					if err != nil {
						logger.Error("'CompletePoll' return error", slog.String("error", err.Error()))
					}
				})
				wg.Done()
			}()

			logger.Info("waiting...")
			wg.Wait()
			logger.Info("load test done...")

			got, err := client.GRPC().GetQueue(context.Background(), &sdkpb.GetQueueRequest{})
			if err != nil {
				logger.Error("error on getting queue", slog.String("err", err.Error()))
			} else {
				logger.Info("get queue return ok", slog.Int64("IdleQueueCount", got.Data.IdleQueueCount), slog.Int64("ActiveQueueCount", got.Data.ActiveMessageCount))
			}
		},
	}
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
