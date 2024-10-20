package http

import (
	"broker/server"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type httpserver struct {
	server http.Server
}

func New(queue server.IBroker) *httpserver {

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

	return &httpserver{
		// queue: queue,
		server: http.Server{
			Addr:    ":2000",
			Handler: r,
		},
	}
}

func (s *httpserver) Start() error { return s.server.ListenAndServe() }
func (s *httpserver) Stop() error  { return s.server.Shutdown(context.Background()) }
