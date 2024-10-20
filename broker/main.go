package main

import (
	"log"
	nethttp "net/http"

	"broker/broker"
	grpcserver "broker/server/grpc"
	httpserver "broker/server/http"
	sdksignal "sdk/signal"
)

func main() {
	log.Printf("hello broker!")

	// setup app's tracer
	shutdownHandler := NewTracer()
	defer shutdownHandler()

	queue := broker.New() // initiate core queue
	queue.Start()         // restore backed-up queues

	httpServer := httpserver.New(queue) // initiate http server
	grpcServer := grpcserver.New(queue) // iniitate grpc server

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
	queue.Stop() // shutdown queue
}
