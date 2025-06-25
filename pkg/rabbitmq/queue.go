package rabbitmq

import (
	"context"
	ampq "github.com/rabbitmq/amqp091-go"
	"github.com/rafian-git/go-backend/pkg/log"
	"time"
)

type Queue interface {
	SetLogger(logger *log.Logger)
	//Connect(uri string) (*ampq.Connection, error)
	ConnectWithRetry(uri string, retries int, delay time.Duration) (*ampq.Connection, error)
	//Consume(ctx context.Context, ch *ampq.Channel, queueName string) (<-chan ampq.Delivery, error)
	ConsumeWithRetry(ctx context.Context, ch *ampq.Channel, queueName string) (<-chan ampq.Delivery, error)
	Publish(ctx context.Context, queueName string, obj []byte) error
	PublishWithRetry(ctx context.Context, queueName string, obj []byte) error
	PublishBulkMessage(ctx context.Context, queueName string, messages []interface{}) error

	ListenDataWithRetry(ctx context.Context, maxRetry int, queueName string, fn ListenDataCallback)
	Listen(ctx context.Context, req *ListenDataReq)
}
