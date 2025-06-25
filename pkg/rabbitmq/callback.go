package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"reflect"
	"time"
)

type ListenDataCallback func(ctx context.Context, requestData interface{}) (interface{}, error)

func (r *RabbitMQ) ListenDataWithRetry(ctx context.Context, maxRetry int, queueName string, fn ListenDataCallback) {
	r.log.Info(ctx, fmt.Sprintf("listening %v  message", queueName))
	var conn *amqp.Connection
	var ch *amqp.Channel
	var msgs <-chan amqp.Delivery
	var err error

	for {
		conn, err = r.ConnectWithRetry(r.config.Url, r.config.Retries, time.Second)
		if err != nil {
			r.log.Error(ctx, fmt.Sprintf("failed to create connection"))
			time.Sleep(5 * time.Second)
			continue
		}

		ch, err = conn.Channel()
		if err != nil {
			r.log.Error(ctx, "failed to connect with channel")
			time.Sleep(5 * time.Second)
			continue
		}

		msgs, err = r.ConsumeWithRetry(ctx, ch, queueName)
		if err != nil {
			r.log.Error(ctx, fmt.Sprintf("failed to consume %s queue", queueName))

			time.Sleep(5 * time.Second)
			continue
		}

		r.log.Info(ctx, fmt.Sprintf("Consuming %s queue", queueName))

		for msg := range msgs {
			r.log.Info(ctx, fmt.Sprintf("received message from %s", queueName), zap.Any("data length", len(msg.Body)))
			if err := r.ProcessMessage(ctx, msg, fn); err != nil {
				if msg.Headers["x-death"] != nil {
					for _, death := range msg.Headers["x-death"].([]interface{}) {
						deathMap := death.(amqp.Table)
						if deathMap["reason"] == "expired" {
							count, ok := deathMap["count"].(int)
							if ok && count < maxRetry {
								_ = msg.Nack(false, false)
							} else {
								//msg.Ack(true)
								r.log.Info(ctx, fmt.Sprintf("maximum retries has been exceeded: %v", msg.MessageId))
							}
							break
						}
					}
				}
			} else {
				msg.Ack(false)
			}
		}

		ch.Close()
		conn.Close()
	}
}

func (r *RabbitMQ) ProcessMessage(ctx context.Context, msg amqp.Delivery, fn ListenDataCallback) error {
	b := bytes.NewBuffer(msg.Body)

	// Create a new instance of the type dynamically
	var requestData interface{}

	// Decode using a dynamic decoding function
	if err := decodeWithReflection(b, &requestData); err != nil {
		errMsg := "Failed to decode message"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", errMsg, err.Error()))
		return err
	}

	_, err := fn(ctx, requestData)
	if err != nil {
		errMsg := "Failed during callback"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", errMsg, err.Error()))
		return err
	}

	return nil
}

// decodeWithReflection decodes the data from the buffer into the provided interface{}
func decodeWithReflection(buffer *bytes.Buffer, target interface{}) error {
	dec := gob.NewDecoder(buffer)

	// Create a pointer to the target type
	targetType := reflect.TypeOf(target)
	targetPtr := reflect.New(targetType).Interface()

	// Decode into the target pointer
	if err := dec.Decode(targetPtr); err != nil {
		return err
	}

	// Copy the decoded value to the original target
	reflect.ValueOf(target).Elem().Set(reflect.ValueOf(targetPtr).Elem())
	return nil
}
