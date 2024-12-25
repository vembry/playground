package http

import (
	"context"
	"fmt"
	"log"
	nethttp "net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/propagation"
)

type server struct {
	server *nethttp.Server
}

func New(addr string, handlers ...*nethttp.ServeMux) *server {
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

// middlewarex is a testing middleware to check request content
func middlewarex(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		handler := otelhttp.NewHandler(
			otelhttp.WithRouteTag(
				r.URL.Path,
				nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
					next.ServeHTTP(w, r)
				}),
			),
			fmt.Sprintf("%s %s", r.Method, r.URL.Path),
			otelhttp.WithPropagators(propagation.TraceContext{}),
		)

		handler.ServeHTTP(w, r)
	})
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
