package handler

import (
	"api/internal/model"
	"context"

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
) *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	s := newHandler(transactionDomain, balanceDomain, addBalanceHandler)

	// create transaction
	r.POST("/transaction", s.CreateTransaction)

	// to create topup
	r.POST("/topup", s.CreateTopup)

	// to get balance
	r.GET("/balance", s.GetBalance)

	// to do health-check
	r.GET("/health", s.HealthCheck)

	return r
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
