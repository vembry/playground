package app

import (
	"app/cmd"

	"github.com/spf13/cobra"
)

// NewCli is to construct clis
func NewCli(server *Server, worker *Worker) *cobra.Command {
	command := &cobra.Command{}
	command.AddCommand(cmd.NewServe(server, worker))
	command.AddCommand(cmd.NewWork(worker))

	return command
}
