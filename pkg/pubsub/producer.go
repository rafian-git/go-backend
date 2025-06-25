package pubsub

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/segmentio/kafka-go"
)

var (
	writer *kafka.Writer
	once   sync.Once
)

func GetSingletonWriter(config *Config) *kafka.Writer {
	once.Do(func() {
		writer = GetWriter(config)
	})
	return writer
}

// pass an existing topic/pass a string to create a topic and send messages to it
func GetWriter(config *Config) *kafka.Writer {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(config.Address),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: config.AllowAutoTopicCreation,
		Compression:            kafka.Snappy,
		BatchSize:              config.BatchSize,
		Async:                  true,
	}

	if config.BatchSize > 0 {
		w.BatchSize = config.BatchSize
	}

	if config.BatchTimeout > 0 {
		w.BatchTimeout = time.Duration(config.BatchTimeout) * time.Millisecond
	}

	return w
}

func PublishMessage(ctx context.Context, writer *kafka.Writer, message kafka.Message) error {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	fmt.Printf("Producer config : TCP-%s ~ topic : %s \n", writer.Addr, message.Topic)
	fmt.Println("kafka message byte size: ", len(message.Value))
	err := writer.WriteMessages(ctx, message)
	if err == nil {
		fmt.Printf("Sent message : data-size %d\n", unsafe.Sizeof(message.Value))
	} else if err == context.Canceled {
		fmt.Println("Context canceled: ", err)
	} else {
		fmt.Println("Error sending message: ", err)
	}

	return err
}

func Publish(config *Config, topic string, message kafka.Message) error {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)
	// ctx, cancel := context.WithCancel(context.Background())

	// go routine for getting signals asynchronously
	// go func() {
	// 	sig := <-signals
	// 	fmt.Printf("Got signal: %s ~ 🛑Stopping Producer \n", sig)
	// 	cancel()
	// }()
	w := GetWriter(config)

	fmt.Printf("Producer config : TCP-%s ~ topic : %s \n", w.Addr, topic)

	defer func() {
		err := w.Close()
		if err != nil {
			fmt.Println("Error closing producer: ", err)
			return
		}
		fmt.Println("Producer closed")
	}()
	fmt.Println("kafka message byte size: ", len(message.Value))
	err := w.WriteMessages(context.Background(), kafka.Message{Topic: topic, Key: []byte(message.Key), Value: []byte(message.Value)})
	if err == nil {
		fmt.Printf("Sent message : data-size %d\n", unsafe.Sizeof(message.Value))
	} else if err == context.Canceled {
		fmt.Println("Context canceled: ", err)
	} else {
		fmt.Println("Error sending message: ", err)
	}

	return err
}

func PublishMessages(ctx context.Context, config *Config, topic string, messages []kafka.Message) error {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	w := GetWriter(config)
	w.Topic = topic
	msg := fmt.Sprintf("Producer config : TCP-%s ~ topic : %s \n", w.Addr, topic)
	fmt.Println(msg)

	defer func() {
		err := w.Close()
		if err != nil {
			msg := fmt.Sprintf("Error closing producer.")
			fmt.Println(msg, err)

			return
		}
		msg := fmt.Sprintf("Producer closed")
		fmt.Println(msg, err)

	}()

	err := w.WriteMessages(ctx, messages...)
	if err == nil {
		msg := fmt.Sprintf("Sent message : data-size %d\n", unsafe.Sizeof(messages))
		fmt.Println(msg)
	} else if errors.Is(err, context.Canceled) {
		msg := fmt.Sprintf("Context canceled.")
		fmt.Println(msg, err)
	} else {
		msg := fmt.Sprintf("Error sending message.")
		fmt.Println(msg, err)
	}

	return err
}
