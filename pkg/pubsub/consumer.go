package pubsub

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"gitlab.techetronventures.com/core/backend/utility"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func Consume(config Config, topic string, group string, c chan kafka.Message) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGKILL)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-signals
		fmt.Printf("Got signal: %s ~ 🛑Stopping Consumer \n", sig)
		cancel()
	}()

	bootstrapServers := strings.Split(utility.GetEnv(utility.BootstrapServers, config.Address), ",")
	topic = utility.GetEnv(utility.Topic, topic)
	group = utility.GetEnv(utility.GroupID, group)
	readerConfig := kafka.ReaderConfig{
		Brokers: bootstrapServers,
		GroupID: group,
		Topic:   topic,
	}

	if config.MaxWait > 0 {
		readerConfig.MaxWait = time.Duration(config.MaxWait) * time.Millisecond
	}

	if config.MaxBytes > 0 {
		readerConfig.MinBytes = 1
		readerConfig.MaxBytes = config.MaxBytes
	}

	r := kafka.NewReader(readerConfig)

	defer func() {
		if err := r.Close(); err != nil {
			fmt.Println("Error closing consumer: ", err)
		}
		fmt.Println("Consumer closed")
	}()

	for {

		m, err := r.ReadMessage(ctx)
		if err != nil {
			fmt.Printf("Error reading message: %v", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		c <- m
	}

	//for msg := range c {
	//	go func(m kafka.Message) {
	//		fmt.Printf("Processed message from %s-%d [%d]: %s = data-size %d\n",
	//			m.Topic, m.Partition, m.Offset, string(m.Key), unsafe.Sizeof(m.Value))
	//	}(msg)
	//}
}
