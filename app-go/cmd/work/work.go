package work

import (
	"app/common"
	"app/internal/worker"
	"log"

	"github.com/spf13/cobra"
)

func New(consumers ...worker.IConsumer) *cobra.Command {
	return &cobra.Command{
		Use:   "work",
		Short: "start worker",
		Long:  "start worker",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("starting worker...")
			for i := range consumers {
				consumers[i].Start()
			}

			common.WatchForExitSignal()
			log.Printf("shutting down worker...")

			for i := range consumers {
				consumers[i].Stop()
			}
		},
	}
}
