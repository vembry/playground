package cmd

import (
	"app/common"
	"log"

	"github.com/spf13/cobra"
)

type IServer interface {
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
			log.Printf("starting server...")
			// h := newHandler(balanceDomain)
			// s := newServer(metric, h)

			// start servers
			for i := range servers {
				servers[i].Start()
			}

			// await
			common.WatchForExitSignal()
			log.Printf("shutting down server...")

			// stop servers gracefully
			for i := range servers {
				servers[i].Stop()
			}
		},
	}
}
