package main

import (
	"load-test/cmd/app"
	"load-test/cmd/broker"
	"log"

	"github.com/spf13/cobra"
)

func main() {
	// setup tracer
	shutdownHandler := newTelemetry()
	defer shutdownHandler()

	// setup logger
	logger := newLogger()

	// setup cli
	cli := cobra.Command{}
	cli.AddCommand(
		broker.New(logger),
		app.New(logger),
	)

	if err := cli.Execute(); err != nil {
		log.Fatalf("found error on executing app's cli. err=%v", err)
	}
}
