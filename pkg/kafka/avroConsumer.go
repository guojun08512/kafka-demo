package kafka

import (
	"context"
	"encoding/binary"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/linkedin/goavro/v2"
)

type avroConsumer struct {
	Consumer             sarama.ConsumerGroup
	Topic                string
	SchemaRegistryClient *CachedSchemaRegistryClient
	handler              *groupConsumerHandler
}

type groupConsumerHandler struct {
	ready chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (handler *groupConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(handler.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (handler *groupConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (handler *groupConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}
	return nil
}

type Message struct {
	SchemaId  int
	Topic     string
	Partition int32
	Offset    int64
	Key       string
	Value     string
}

func NewConsumerConfig() (conf *sarama.Config) {
	conf = sarama.NewConfig()
	conf.Version = sarama.MaxVersion
	conf.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	conf.Consumer.Offsets.AutoCommit.Enable = true
	conf.Consumer.Offsets.Initial = sarama.OffsetOldest
	conf.Consumer.Fetch.Min = 1
	conf.Consumer.Fetch.Default = 1024
	conf.Consumer.MaxWaitTime = time.Millisecond * 100

	// FIXME: For consumer it's per-partition/channel value. It's default value is 256.
	//  May cause huge memory usage (partition_count*buffer_size*message_size).
	conf.ChannelBufferSize = 10
	return
}

// avroConsumer is a basic consumer to interact with schema registry, avro and kafka
func NewAvroConsumer(kafkaServers []string, schemaRegistryServers []string,
	topic string, groupId string, callbacks *groupConsumerHandler) (*avroConsumer, error) {
	// init (custom) config, enable errors and notifications
	config := NewConsumerConfig()
	config.Consumer.Return.Errors = true
	//read from beginning at the first time
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup(kafkaServers, groupId, config)
	if err != nil {
		return nil, err
	}

	schemaRegistryClient := NewCachedSchemaRegistryClient(schemaRegistryServers)
	if callbacks == nil {
		callbacks = &groupConsumerHandler{
			ready: make(chan bool),
		}
	}
	return &avroConsumer{
		consumer,
		topic,
		schemaRegistryClient,
		callbacks,
	}, nil
}

//GetSchemaId get schema id from schema-registry service
func (ac *avroConsumer) GetSchema(id int) (*goavro.Codec, error) {
	codec, err := ac.SchemaRegistryClient.GetSchema(id)
	if err != nil {
		return nil, err
	}
	return codec, nil
}

func (ac *avroConsumer) Consume() {
	// trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the handler session will need to be
			// recreated to get the new claims
			if err := ac.Consumer.Consume(ctx, []string{ac.Topic}, ac.handler); err != nil {
				log.WithError(err).Warn("kafka consumer error")
			}
			if ctx.Err() != nil {
				return
			}
			log.Warnf("kafka consumer session closed, topics=(%s), need reconnect", ac.Topic)
			ac.handler.ready = make(chan bool)
		}
	}()
	log.Println("Sarama consumer up and running!...")
	<-ac.handler.ready
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-sigterm:
		log.Println("terminating: via signal")
	}
	cancel()
	wg.Wait()
	if err := ac.Consumer.Close(); err != nil {
		log.Panicf("Error closing client: %v", err)
	}
}

func (ac *avroConsumer) ProcessAvroMsg(m *sarama.ConsumerMessage) (Message, error) {
	schemaId := binary.BigEndian.Uint32(m.Value[1:5])
	codec, err := ac.GetSchema(int(schemaId))
	if err != nil {
		return Message{}, err
	}
	// Convert binary Avro data back to native Go form
	native, _, err := codec.NativeFromBinary(m.Value[5:])
	if err != nil {
		return Message{}, err
	}

	// Convert native Go form to textual Avro data
	textual, err := codec.TextualFromNative(nil, native)

	if err != nil {
		return Message{}, err
	}
	msg := Message{int(schemaId), m.Topic, m.Partition, m.Offset, string(m.Key), string(textual)}
	return msg, nil
}

func (ac *avroConsumer) Close() {
	ac.Consumer.Close()
}
