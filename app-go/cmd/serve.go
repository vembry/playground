package cmd

import (
	"app/common"
	"log"

	"github.com/spf13/cobra"
)

// serverProvider contain the spec of a 'server'
type serverProvider interface {
	Start() error
	PostStartCallback()
	PostShutdownCallback()
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

			// run post-start callback
			server.PostStartCallback()

			common.WatchForExitSignal()

			log.Print("* Shutting down the server...")
			if err := server.Shutdown(); err != nil {
				log.Printf("found error on shutting down server. err=%v", err)
			}

			// run post-shutdown callback
			server.PostShutdownCallback()
		},
	}
}
