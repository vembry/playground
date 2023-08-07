package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// serverProvider contain the spec of a 'server'
type serverProvider interface {
	Start() error
	Shutdown() error
	GetAddress() string
}

// NewServe is to initiate cli command of 'serve'
func NewServe(server serverProvider, worker workerProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run app's http server",
		Run: func(cmd *cobra.Command, args []string) {
			go func() {
				log.Printf("* Starting the server at %s...", server.GetAddress())
				server.Start()
			}()

			// had to do this, otherwise worker's handler
			// wont be able to enqueue task to worker
			if worker != nil {
				worker.ConnectToQueue()
			}

			watchForExitSignal()

			log.Print("* Shutting down the server...")
			if err := server.Shutdown(); err != nil {
				log.Printf("found error on shutting down server. err=%v", err)
			}

			// had to do this, because of prior worker.ConnectToQueue
			if worker != nil {
				if err := worker.DisconnectFromQueue(); err != nil {
					log.Printf("found error on disconnecting from queue. err=%v", err)
				}
			}
		},
	}
}
