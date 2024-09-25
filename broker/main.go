package main

import (
	"log"
	"net/http"

	brokergrpc "broker/grpc"
	brokerhttp "broker/http"
	sdksignal "sdk/signal"
)

func main() {
	log.Printf("hello broker!")

	queue := newQueue() // initiate core queue
	queue.restore()     // restore backed-up queues

	httpServer := brokerhttp.NewHttpServer(queue) // initiate http server
	grpcServer := brokergrpc.NewGrpcServer(queue) // iniitate grpc server

	// http server
	go func() {
		if err := httpServer.Start(); err != http.ErrServerClosed {
			log.Fatalf("found error on starting http server. err=%v", err)
		}
	}()

	// grpc server
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("found error on starting grpc server. err=%v", err)
		}
	}()

	sdksignal.WatchForExitSignal()
	log.Println("shutting down...")

	httpServer.Stop()
	grpcServer.Stop()

	queue.shutdown() // shutdown queue
}
