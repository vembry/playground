package cmd

import (
	"app-go/common"
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/spf13/cobra"
)

// WorkerHandler contains the spec of a worker handler. This is for consuming-worker
type WorkerHandler interface {
	Type() string
	Perform(context.Context, *asynq.Task) error
}

// workerProvider contains the spec of worker
type workerProvider interface {
	Start() error
	Shutdown() error
	PostStartCallback()
	PostShutdownCallback()
}

// NewWork is to initiate cli command of 'work'
func NewWork(worker workerProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "work",
		Short: "Run app's worker",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("* Starting the worker...")

			if err := worker.Start(); err != nil {
				log.Fatalf("found error on starting worker. err=%v", err)
			}

			// run post-start callback
			worker.PostStartCallback()

			common.WatchForExitSignal()

			log.Print("* Shutting down the worker...")
			if err := worker.Shutdown(); err != nil {
				log.Printf("found error on shutting down server. err=%v", err)
			}

			// run post-shutdown callback
			worker.PostShutdownCallback()
		},
	}
}
