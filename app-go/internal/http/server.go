package http

import (
	"app/internal/domain"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type server struct {
	httpserver *http.Server
}

func (s *server) Name() string {
	return "http"
}

type IMetric interface {
	RecordInbound(route string, method string, statusCode string, duration time.Duration)
}

func NewServer(
	httpaddress string,
	metric IMetric,
	balanceDomain domain.IBalance,
) *server {

	h := newHandler(balanceDomain)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(inboundMetric(metric))

	r.POST("balance/open", h.Open)
	r.GET("balance/:balance_id", h.Get)
	r.POST("balance/:balance_id/deposit", h.Deposit)
	r.POST("balance/:balance_id/withdraw", h.Withdraw)
	r.POST("balance/:balance_id/transfer", h.Transfer)

	// to do health-check
	r.GET("/health", h.HealthCheck)

	mux := http.NewServeMux()
	mux.Handle("/", r.Handler())

	return &server{
		httpserver: &http.Server{
			Addr:    httpaddress,
			Handler: mux,
		},
	}
}

func (s *server) Start() {
	go func() {
		if err := s.httpserver.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("found error on starting http server. err=%v", err)
		}
	}()
}
func (s *server) Stop() {
	// context for stop timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// stop server
	err := s.httpserver.Shutdown(ctx)
	if err != nil {
		log.Printf("found error on stopping http server. err=%v", err)
	}
}
