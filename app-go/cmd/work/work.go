package work

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "work",
		Short: "start worker",
		Long:  "start worker",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
}
