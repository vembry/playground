package cmd

import (
	"log"

	sdksignal "sdk/signal"

	"github.com/spf13/cobra"
)

type IWorker interface {
	Name() string
	Start()
	Stop()
}

func NewWork(metricServer IServer, workers ...IWorker) *cobra.Command {
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

			metricServer.Start()

			sdksignal.WatchForExitSignal()
			log.Printf("shutting down worker...")

			for i := range workers {
				log.Printf("shutting down %s worker", workers[i].Name())
				workers[i].Stop()
			}
			metricServer.Stop()
		},
	}
}
