package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/rand/v2"
	sdkpb "sdk/pb"
	"sdk/tester"

	"github.com/segmentio/ksuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	logger := slog.New(slog.Default().Handler())

	// setup load tester
	t := tester.New(
		tester.Config{
			Type:                  tester.LoadType_Count,
			Logger:                logger,
			ConcurrentWorkerCount: 1,
			Counter:               1000,
		},
	)

	// setup parameter for test
	params := []string{}
	for i := range 5 {
		params = append(params, fmt.Sprintf("queue_%d", i))
	}

	// setup grpc connection
	conn, err := grpc.NewClient("localhost:4000", grpc.WithTransportCredentials(insecure.NewCredentials())) // Use the server's address and port
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := sdkpb.NewBrokerClient(conn)

	// run load tester
	t.Do(func(ctx context.Context, logger *slog.Logger) {
		// choose balance id
		i := randRange(0, len(params)-1)
		queueName := params[i]

		_, err := client.Enqueue(ctx, &sdkpb.EnqueueRequest{
			QueueName: queueName,
			Payload:   ksuid.New().String(),
		})

		// break away when theres error
		if err != nil {
			log.Panicf("load test panics. err=%v", err)
		}
	})
}

func randRange(min, max int) int {
	return rand.IntN(max-min) + min
}
