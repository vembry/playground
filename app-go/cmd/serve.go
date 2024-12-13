package cmd

import (
	"context"
	"log"
	sdksignal "sdk/signal"
	"time"

	"github.com/spf13/cobra"
)

type IServer interface {
	Name() string
	Start()
	Stop(context.Context)
}

func NewServe(
	servers ...IServer,
) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "start server",
		Long:  "start server long",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("starting servers...")

			// start servers
			for i := range servers {
				log.Printf("starting %s server", servers[i].Name())
				servers[i].Start()
			}

			// await
			sdksignal.WatchForExitSignal()
			log.Printf("shutting down server...")

			// stop servers gracefully
			for i := range servers {
				// context for stop timeout
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				log.Printf("shutting down %s server", servers[i].Name())

				servers[i].Stop(ctx)
			}
		},
	}
}
