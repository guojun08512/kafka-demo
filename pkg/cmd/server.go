package cmd

import (
	"github.com/spf13/cobra"
	"keyayun.com/seal-kafka-runner/pkg/kafka"
)

func startUp() error {
	kafka.StartUpConsumer()
	return nil
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "runner-server serve",
	RunE: func(cmd *cobra.Command, Args []string) error {
		return startUp()
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
}
