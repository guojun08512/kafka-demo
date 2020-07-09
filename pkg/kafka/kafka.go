package kafka

import (
	"keyayun.com/seal-kafka-runner/pkg/config"
	"keyayun.com/seal-kafka-runner/pkg/logger"
)

var (
	conf  = config.Config
	log   = logger.WithNamespace("kafka")
	group = "seal-runner-kafka"
)

func StartUpConsumer() error {
	brokers := conf.GetStringSlice("kafka.brokers")
	schemaRegistries := conf.GetStringSlice("kafka.schemaRegistries")
	client, err := NewAvroConsumer(brokers, schemaRegistries, "test", group, nil)
	if err != nil {
		log.Panicf("Error creating consumer group client: %v", err)
		return err
	}
	client.Consume()
	return nil
}

func NewSyncProducer() (*AvroProducer, error) {
	brokers := conf.GetStringSlice("kafka.brokers")
	schemaRegistries := conf.GetStringSlice("kafka.schemaRegistries")
	return NewAvroProducer(brokers, schemaRegistries)
}
