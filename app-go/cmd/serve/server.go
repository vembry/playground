package serve

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type server struct {
	httpserver *http.Server
}

func newServer(h *handler) *server {
	r := gin.Default()

	r.POST("balance/open", h.Open)
	r.GET("balance/:balance_id", h.Get)
	r.POST("balance/:balance_id/deposit", h.Deposit)
	r.POST("balance/:balance_id/withdraw", h.Withdraw)
	r.POST("balance/:balance_id/transfer", h.Transfer)

	mux := http.NewServeMux()
	mux.Handle("/", r.Handler())

	return &server{
		httpserver: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		},
	}
}

func (s *server) Start() {
	log.Print("starting server...")

	log.Print("starting http server...")
	go func() {
		if err := s.httpserver.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("found error on starting http server. err=%v", err)
		}
	}()
}
func (s *server) Stop() {
	log.Print("stopping server...")

	// context for stop timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// stop server
	err := s.httpserver.Shutdown(ctx)
	if err != nil {
		log.Printf("found error on stopping http server. err=%v", err)
	}
}
