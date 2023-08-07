package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

type ServerProvider interface {
	Start() error
	GracefulStop() error
	GetAddress() string
}

// NewServe is to initiate cli command of 'serve'
func NewServe(server ServerProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run http server",
		Run: func(cmd *cobra.Command, args []string) {
			go func() {
				log.Printf("* Starting the server at %s...", server.GetAddress())
				server.Start()
			}()

			watchForExitSignal()

			log.Print("* Shutting down the server...")
			if err := server.GracefulStop(); err != nil {
				log.Printf("found error on shutting down server. err=%v", err)
			}

		},
	}
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
