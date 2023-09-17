package app

import (
	"context"
	"net/http"
	"time"
)

// Server contain the config of app-server
type Server struct {
	httpAddress          string
	server               *http.Server
	postStartCallback    func()
	postShutdownCallback func()
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

// Shutdown is to shutdown server gracefully
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

// GetAddress is to get server address
func (s *Server) GetAddress() string {
	return s.httpAddress
}

// WithPostStartCallback is to inject callback on post-start
func (s *Server) WithPostStartCallback(callback func()) {
	s.postStartCallback = callback
}

// PostStartCallback executes callback on post-start
func (s *Server) PostStartCallback() {
	// return s.postStartCallback
	if s.postStartCallback != nil {
		s.postStartCallback()
	}
}

// WithPostShutdownCallback is to inject callback on post-shutdown
func (s *Server) WithPostShutdownCallback(callback func()) {
	s.postShutdownCallback = callback
}

// PostShutdownCallback executes callback on post-shutdown
func (s *Server) PostShutdownCallback() {
	// do post-shutdown
	if s.postShutdownCallback != nil {
		s.postShutdownCallback()
	}
}
