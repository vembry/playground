package http

import (
	"broker/model"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type server struct {
	server http.Server
}

type IQueue interface {
	Get() model.QueueData
	Enqueue(payload model.EnqueuePayload) error
	Poll(queueName string) (*model.ActiveQueue, error)
	CompletePoll(queueId ksuid.KSUID) error
}

func NewHttpServer(queue IQueue) *server {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	h := handler{
		queue: queue,
	}

	// queue handler
	queueGroup := r.Group("/queue")
	queueGroup.GET("", h.get)
	queueGroup.POST("/enqueue", h.enqueue)
	queueGroup.GET("/poll/:queue_name", h.poll)
	queueGroup.POST("/poll/:queue_id/complete", h.completePoll)

	return &server{
		// queue: queue,
		server: http.Server{
			Addr:    ":2000",
			Handler: r,
		},
	}
}

func (s *server) Start() error { return s.server.ListenAndServe() }
func (s *server) Stop() error  { return s.server.Shutdown(context.Background()) }
