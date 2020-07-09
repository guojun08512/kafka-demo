package cmd

import (
	"github.com/Shopify/sarama"
	"github.com/spf13/cobra"
	"keyayun.com/seal-kafka-runner/pkg/errors"
	"keyayun.com/seal-kafka-runner/pkg/kafka"
)

var producer = &cobra.Command{
	Use:   "producer",
	Short: "producer",
	RunE: func(cmd *cobra.Command, args []string) error {
		return producerTopic(cmd, args)
	},
}

func producerTopic(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.InvalidArg()
	}
	key := args[0]
	val := args[1]
	pro, err := kafka.NewSyncProducer()
	if err != nil {
		return err
	}
	kafkaKey, err := sarama.StringEncoder(key).Encode()
	if err != nil {
		return err
	}
	kafkaVal, err := sarama.StringEncoder(val).Encode()
	if err != nil {
		return err
	}
	return pro.Add("test", "string", kafkaKey, kafkaVal)
}

var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "runner-server kafka",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	kafkaCmd.AddCommand(producer)
	RootCmd.AddCommand(kafkaCmd)
}
