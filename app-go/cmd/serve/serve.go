package serve

import (
	"app/common"
	"app/internal/app"
	"app/internal/domain"
	"log"

	"github.com/spf13/cobra"
)

func New(
	cfg *app.EnvConfig,
	balanceDomain domain.IBalance,
) *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "start server",
		Long:  "start server long",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("starting server...")
			h := newHandler(balanceDomain)
			s := newServer(h)

			// start server
			s.Start()

			// await
			common.WatchForExitSignal()
			log.Printf("shutting down server...")

			// stop server gracefully
			s.Stop()
		},
	}
}
