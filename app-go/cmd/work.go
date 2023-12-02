package cmd

import (
	"app/common"
	"log"

	"github.com/spf13/cobra"
)

type IWorker interface {
	Name() string
	Start()
	Stop()
}

func NewWork(workers ...IWorker) *cobra.Command {
	return &cobra.Command{
		Use:   "work",
		Short: "start worker",
		Long:  "start worker",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("starting worker...")
			for i := range workers {
				log.Printf("starting %s worker", workers[i].Name())
				workers[i].Start()
			}

			common.WatchForExitSignal()
			log.Printf("shutting down worker...")

			for i := range workers {
				log.Printf("shutting down %s worker", workers[i].Name())
				workers[i].Stop()
			}
		},
	}
}
