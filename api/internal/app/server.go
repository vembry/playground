package app

import (
	"context"
	"net/http"
	"time"
)

// Server contain the config of app-server
type Server struct {
	httpAddress string
	server      *http.Server
}

// NewServer is to setup app-server
func NewServer(cfg *EnvConfig, handler http.Handler) *Server {

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	return &Server{
		httpAddress: cfg.HttpAddress,

		server: &http.Server{
			Addr:    cfg.HttpAddress,
			Handler: mux,
		},
	}
}

// Start is to start server
func (s *Server) Start() error {
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// GracefulStop is to shutdown server gracefully
func (s *Server) GracefulStop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	s.server.Shutdown(ctx)

	return s.server.Shutdown(context.Background())
}

// GetAddress is to get server address
func (s *Server) GetAddress() string {
	return s.httpAddress
}