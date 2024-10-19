package main

import (
	"load-test/cmd/app"
	"load-test/cmd/broker"

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
}
