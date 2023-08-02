package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

// CreateTransactionRequest contains fields submitted via create-transaction api
type CreateTransactionRequest struct {
	Amount      float64 `json:"amount"`
	Type        string  `json:"type"` // credit or debit
	Description string  `json:"description"`
}

// GenericResponse contains fields returned to api requester
type GenericResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// main is to start the server
func main() {
	r := gin.Default()

	// create transaction
	r.POST("/transaction", func(ctx *gin.Context) {

		var trxRequest CreateTransactionRequest
		if err := ctx.Bind(&trxRequest); err != nil {
			// when error, return 4xx
			ctx.JSON(400, GenericResponse{
				Status:  false,
				Message: fmt.Errorf("found error. err=%w", err).Error(),
			})
			return
		}

		ctx.JSON(200, map[string]interface{}{
			"status":  true,
			"payload": trxRequest,
		})
	})

	log.Printf("starting http server...")

	go func() {
		if err := r.Run("0.0.0.0:8080"); err != nil {
			log.Fatalf("gin stopped running. err=%v", err)
		}
	}()

	// awaits for interrupt signals
	watchForExitSignal()

	log.Printf("stopping http server...")

	// do shutdown handling
	// ...

	log.Printf("server stopped")
}

// watchForExitSignal is to awaits incoming interrupt signal
// sent to the service
func watchForExitSignal() os.Signal {
	ch := make(chan os.Signal, 4)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	)

	return <-ch
}
