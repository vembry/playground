package main

import (
	"log"
	nethttp "net/http"

	"broker/grpc"
	"broker/http"
	"broker/queue"
	sdksignal "sdk/signal"
)

func main() {
	log.Printf("hello broker!")

	// setup app's tracer
	shutdownHandler := NewTracer()
	defer shutdownHandler()

	queue := queue.New()   // initiate core queue
	queue.Start()          // restore backed-up queues
	defer queue.Shutdown() // shutdown queue

	httpServer := http.NewServer(queue) // initiate http server
	grpcServer := grpc.NewServer(queue) // iniitate grpc server

	// http server
	go func() {
		if err := httpServer.Start(); err != nethttp.ErrServerClosed {
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
}
