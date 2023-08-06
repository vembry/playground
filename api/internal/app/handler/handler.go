package handler

import (
	"api/internal/model"
	"context"

	"github.com/gin-gonic/gin"
)

// transactionProvider is the spec of transaction-domain instance
type transactionProvider interface {
	Create(ctx context.Context, in *model.CreateTransaction) (*model.CommonResponse, error)
}

// ledgerProvider is the spec of ledger-domain instance
type ledgerProvider interface {
}

// GenericResponse contains fields returned to api requester
type GenericResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// NewHttpHandler is to setup http handler
func NewHttpHandler(transactionDomain transactionProvider, ledgerDomain ledgerProvider) *gin.Engine {

	r := gin.Default()
	s := newServer(transactionDomain)

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

// server contain the server instance to handle incoming http request
type server struct {
	transactionDomain transactionProvider
}

// newServer is to initiate server
func newServer(transactionDomain transactionProvider) *server {
	return &server{
		transactionDomain: transactionDomain,
	}
}
