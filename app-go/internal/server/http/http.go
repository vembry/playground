package http

import (
	"context"
	"log"
	"net/http"
	nethttp "net/http"
)

type server struct {
	server *nethttp.Server
}

func New(addr string, handlers ...http.Handler) *server {
	mux := nethttp.NewServeMux()

	// register all handler into the server mux
	for _, handler := range handlers {
		mux.Handle("/", handler)
	}

	// manually inject health check endpoints
	mux.HandleFunc("GET /health", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		w.WriteHeader(nethttp.StatusOK)
	})

	return &server{
		server: &nethttp.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *server) Name() string {
	return "core"
}

func (s *server) Start() {
	go func() {
		if err := s.server.ListenAndServe(); err != nethttp.ErrServerClosed {
			log.Fatalf("found error on starting http server. err=%v", err)
		}
	}()
}

// Stop will require caller to pass context with timeout
func (s *server) Stop(ctx context.Context) {
	// stop server
	err := s.server.Shutdown(ctx)
	if err != nil {
		log.Printf("found error on stopping http server. err=%v", err)
	}
}
