package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"strings"
	"time"
)

func (k *Kafka) getSingletonWriter() *kafka.Writer {
	once.Do(func() {
		writer = createWriter(k.config)
	})
	return writer
}

// ConsumeMessage continuously reads messages and writes them to the provided channel.
func (k *Kafka) consumeMessage(ctx context.Context, topic string, group string, c chan kafka.Message) {

	bootstrapServers := strings.Split(k.config.Address, ",")

	readerConfig := kafka.ReaderConfig{
		Brokers: bootstrapServers,
		GroupID: group,
		Topic:   topic,
		MaxWait: 500 * time.Millisecond,
	}

	reader := kafka.NewReader(readerConfig)
	fmt.Println(readerConfig)
	defer func() {
		err := reader.Close()
		if err != nil {
			fmt.Println("Error closing consumer:", err)
		} else {
			fmt.Println("Consumer closed")
		}
	}()

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			// Handle message read error (can add retry logic if needed)
			fmt.Println("Error reading message:", err)
			time.Sleep(1 * time.Second) // Retry with some delay
			continue
		}

		// Write the message to the channel
		c <- msg

		// You can log the message or process it as needed
		fmt.Printf("Received message from %s [%d]: %s = data-size %d\n",
			msg.Topic, msg.Partition, string(msg.Key), len(msg.Value))
	}
}

// startConsumer handles signal notifications and invokes message consumption.
func (k *Kafka) startConsumer(topic, group string, c chan kafka.Message) {
	ctx, cancel := createContextWithSignalHandling()
	defer cancel()

	// Consume messages
	k.consumeMessage(ctx, topic, group, c)
}
