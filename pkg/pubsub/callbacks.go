package pubsub

import (
	"context"
	"fmt"
	"github.com/rafian-git/go-backend/pkg/log"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// callback function type for example topic
type ExampleCallbackFunc func(str string)

// for each topic create a callback function here and call ListenExampleTopic(callbackFunc)
// in a separate go routine from main.go in your service to consume messages from the topic
func ListenExampleTopic(config *Config, fn ExampleCallbackFunc) {
	chanData := make(chan kafka.Message)
	go Consume(*config, "my-topic", "my-grp", chanData)
	for msg := range chanData {
		fmt.Println("printing frm backend - ", string(msg.Value))
		fn(string(msg.Value))
	}
}

type ListenCallback func(ctx context.Context, request []byte)

func Listen(config *Config, log *log.Logger, topic, group string, callback ListenCallback) {
	chanData := make(chan kafka.Message)
	log.Info(context.Background(), fmt.Sprintf("Establishing consumer connection: %s %s", topic, group))
	go Consume(*config, topic, group, chanData)
	for msg := range chanData {
		log.Info(context.Background(), topic, zap.Any("byte-length", len(msg.Value)))
		callback(context.Background(), msg.Value)
	}
}
