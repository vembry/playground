package cmd

import (
	"app/internal/worker/rabbit"
	"log"

	"github.com/spf13/cobra"
)

func NewDummy() *cobra.Command {
	return &cobra.Command{
		Use:   "dumb",
		Short: "start dummy",
		Long:  "start dummy long",
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("starting dummy...")
			rabbit.New("")
		},
	}
}
