package main

import (
	"context"
	"log"
	"net/http"

	sdksignal "sdk/signal"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Printf("hello broker!")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	queue := newQueue()
	queue.restore() // restore backed-up queues

	// queue handler
	queueGroup := r.Group("/queue")
	queueGroup.GET("", queue.get)
	queueGroup.POST("/enqueue", queue.enqueue)
	queueGroup.GET("/poll/:queue_name", queue.poll)
	queueGroup.POST("/poll/:queue_id/complete", queue.completePoll)

	httpserver := http.Server{
		Addr:    ":2000",
		Handler: r,
	}

	go func() {
		if err := httpserver.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("found error on starting http server. err=%v", err)
		}
	}()

	sdksignal.WatchForExitSignal()
	log.Println("shutting down...")

	httpserver.Shutdown(context.Background())

	queue.shutdown()
}
