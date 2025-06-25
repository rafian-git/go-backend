package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gitlab.techetronventures.com/core/backend/pkg/log"
	"go.uber.org/zap"
	"sync"
)

var (
	writer *kafka.Writer
	once   sync.Once
)

// Kafka struct holds the Kafka configuration and logger for use across services.
type Kafka struct {
	config *Config
	logger *log.Logger
}

// New returns a new Kafka instance with the provided config and log instance.
func New(config *Config, logger *log.Logger) *Kafka {
	return &Kafka{
		config: config,
		logger: logger.Named("kafka-"),
	}
}

// Publish sends a message to the configured topic in the Kafka instance.
func (k *Kafka) Publish(topic string, message kafka.Message) error {
	// Create a context with signal handling for graceful shutdown
	ctx, cancel := createContextWithSignalHandling()
	defer cancel()

	// Set the topic if not already set in the message
	if message.Topic == "" {
		message.Topic = topic
	}

	return k.publishMessage(ctx, message)
}

//// PublishAsync sends a Kafka message asynchronously.
//func (k *Kafka) PublishAsync(config *Config, topic string, message kafka.Message) {
//	go func() {
//		if err := k.Publish(topic, message); err != nil {
//			fmt.Println("Error publishing asynchronously:", err)
//		}
//	}()
//}

type ListenCallback func(ctx context.Context, request []byte) error

func (k *Kafka) Subscribe(topic, group string, callback ListenCallback) {
	k.logger = log.New().Named("kafka-listener")
	ctx := context.Background()
	chanData := make(chan kafka.Message)

	k.logger.Info(ctx, fmt.Sprintf("Establishing consumer connection: %s %s", topic, group))

	go k.startConsumer(topic, group, chanData)
	for msg := range chanData {
		k.logger.Info(ctx, topic, zap.Any("byte-length", len(msg.Value)))
		err := callback(ctx, msg.Value)
		if err != nil {
			k.logger.Error(ctx, fmt.Sprintf("Error while consuming kafka message: %v", err.Error()))
			k.Publish(topic+"-dlq", kafka.Message{Value: msg.Value, Key: msg.Key})
		}

	}
}
