package handler

import (
	"app/internal/model"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

// transactionProvider is the spec of transaction-domain instance
type transactionProvider interface {
	Create(ctx context.Context, in *model.CreateTransaction) error
}

// balanceProvider is the spec of balance-domain instance
type balanceProvider interface {
	Add(ctx context.Context, in *model.AddBalanceParam) error
	Get(ctx context.Context, userId ksuid.KSUID) (*model.Balance, error)
}

// addBalanceHandlerProvider contain spec for add-balance handler
type addBalanceHandlerProvider interface {
	Enqueue(ctx context.Context, in *model.AddBalanceParam) error
}

// requestMetricProvider is the spec for request metric aggregator(?)
type requestMetricProvider interface {
	GinRequest(*gin.Context)
}

// GenericResponse contains fields returned to api requester
type GenericResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// NewHttpHandler is to setup http handler
func NewHttpHandler(
	transactionDomain transactionProvider,
	balanceDomain balanceProvider,
	addBalanceHandler addBalanceHandlerProvider,
	requestMetric requestMetricProvider,
) http.Handler {

	gin.SetMode(gin.ReleaseMode)
	s := newHandler(transactionDomain, balanceDomain, addBalanceHandler)

	r := gin.Default()

	r.Use(requestMetric.GinRequest)

	// create transaction
	r.POST("/transaction", s.CreateTransaction)

	// balance group
	balanceGroup := r.Group("/balance")
	balanceGroup.GET("", s.GetBalance)
	balanceGroup.POST("/add", s.AddBalance)

	// to do health-check
	r.GET("/health", s.HealthCheck)

	return r.Handler()
}

// handler contain the server instance to handle incoming http request
type handler struct {
	transactionDomain transactionProvider
	balanceDomain     balanceProvider
	addBalanceHandler addBalanceHandlerProvider
}

// newHandler is to setup handler instance
func newHandler(
	transactionDomain transactionProvider,
	balanceDomain balanceProvider,
	addBalanceHandler addBalanceHandlerProvider,
) *handler {
	return &handler{
		transactionDomain: transactionDomain,
		balanceDomain:     balanceDomain,
		addBalanceHandler: addBalanceHandler,
	}
}
