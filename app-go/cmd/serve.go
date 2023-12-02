package cmd

import (
	"app/common"
	"log"

	"github.com/spf13/cobra"
)

type IServer interface {
	Name() string
	Start()
	Stop()
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
			common.WatchForExitSignal()
			log.Printf("shutting down server...")

			// stop servers gracefully
			for i := range servers {
				log.Printf("shutting down %s server", servers[i].Name())
				servers[i].Stop()
			}
		},
	}
}
