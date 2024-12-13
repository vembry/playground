package cmd

import (
	"fmt"
	"log"
	"strings"

	sdksignal "sdk/signal"

	"github.com/spf13/cobra"
)

type IWorker interface {
	Name() string
	Start()
	Stop()
}

func NewWork(workers ...IWorker) *cobra.Command {
	// setup workers mapper for selection purposes
	availableWorkerMap := map[string]IWorker{}
	availableWorkers := []string{}
	for _, worker := range workers {
		availableWorkerMap[worker.Name()] = worker
		availableWorkers = append(availableWorkers, worker.Name())
	}

	selectedWorkers := []string{}
	c := &cobra.Command{
		Use:   "work",
		Short: "start worker",
		Long:  "start worker",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("starting worker...")

			// when 'selectedWorkers' is empy,
			// default to run everything
			if len(selectedWorkers) == 0 {
				selectedWorkers = availableWorkers
			}

			// store selected worker for shutdown later
			activeWorker := []IWorker{}

			// iterate to run ONLY selected worker
			for _, selectedWorker := range selectedWorkers {
				val := availableWorkerMap[selectedWorker]
				log.Printf("starting '%s' worker", val.Name())
				val.Start()

				activeWorker = append(activeWorker, val)
			}

			sdksignal.WatchForExitSignal()
			log.Printf("shutting down '%d' worker(s)...", len(activeWorker))

			// iteratively shutdown worker
			for _, worker := range activeWorker {
				log.Printf("shutting down '%s' worker", worker.Name())
				worker.Stop()
			}
		},
	}

	// cli flags
	c.Flags().StringArrayVarP(&selectedWorkers, "worker", "w", []string{}, fmt.Sprintf("determine worker(s) to be run. available worker=%s", strings.Join(availableWorkers, ", ")))

	return c
}
