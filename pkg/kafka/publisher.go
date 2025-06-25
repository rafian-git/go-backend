package kafka

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"
)

// CreateWriter initializes and returns a Kafka writer based on the given config.
func createWriter(config *Config) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(config.Address),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: config.AllowAutoTopicCreation,
		Compression:            kafka.Snappy,
		Async:                  config.Async,
	}
}

// publishMessage sends a message using the Kafka writer in the Kafka instance.
func (k *Kafka) publishMessage(ctx context.Context, message kafka.Message) error {
	w := k.getSingletonWriter()
	k.logger.Info(ctx, fmt.Sprintf("Producer config: TCP-%s ~ topic: %s\n", w.Addr, message.Topic))

	k.logger.Info(ctx, fmt.Sprintf("Kafka message byte size: %v\n", len(message.Value)))

	if err := w.WriteMessages(ctx, message); err != nil {
		k.logger.Error(ctx, fmt.Sprintf("Error while sending message %v", err.Error()))
		return fmt.Errorf("error sending message: %w", err)
	}

	k.logger.Info(ctx, "Sent message successfully")
	return nil
}

// createContextWithSignalHandling returns a context and cancel function that
// listens for system signals (SIGINT, SIGKILL) to gracefully shut down.
func createContextWithSignalHandling() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		<-signals
		fmt.Println("Received stop signal")
		cancel()
	}()

	return ctx, cancel
}
