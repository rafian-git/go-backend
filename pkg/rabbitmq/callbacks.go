package rabbitmq

import (
	"context"
	"fmt"
	ampq "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"time"
)

type ListenDataReq struct {
	QueueName     string
	Callback      ListenCallback
	MaxRetries    int64
	RetryInterval int64
}

type ListenCallback func(ctx context.Context, rr []byte) error

func (r *RabbitMQ) Listen(ctx context.Context, req *ListenDataReq) {
	logger := r.log.Named(req.QueueName)
	logger.Info(ctx, fmt.Sprintf("Listen queue : %s", req.QueueName))

	if req.MaxRetries == 0 {
		req.MaxRetries = MaxRetriesConsume
	}

	if req.RetryInterval == 0 {
		req.RetryInterval = DelayFiftySec
	}

	for {
		// TODO: Refactoring candidate
		func(ctx context.Context) {
			var conn *ampq.Connection
			var ch *ampq.Channel
			var msgs <-chan ampq.Delivery
			var err error
			conn, err = r.ConnectWithRetry(r.config.Url, r.config.Retries, time.Second)
			if err != nil {
				logger.Error(ctx, fmt.Sprintf("failed to create connection with %s", req.QueueName))
				time.Sleep(5 * time.Second)
				return
			}
			defer conn.Close()

			ch, err = conn.Channel()
			if err != nil {
				logger.Error(ctx, fmt.Sprintf("failed to open channel with %s", req.QueueName))
				time.Sleep(5 * time.Second)
				return
			}
			defer ch.Close()

			//prefetchCount := 1
			//err = ch.Qos(prefetchCount, 0, false)
			if err != nil {
				logger.Error(ctx, fmt.Sprintf("failed to set prefetch for  %s", req.QueueName))
				return
			}

			//request := &portfolio_pb.PlaceOrderRequest{}
			msgs, err = r.ConsumeRMQ(ctx, ch, req.QueueName, req.RetryInterval)
			if err != nil {
				logger.Error(ctx, fmt.Sprintf("failed to consume %s queue", req.QueueName))
				time.Sleep(5 * time.Second)
				return
			}

			logger.Info(ctx, fmt.Sprintf("Consuming %s queue", req.QueueName))

			for msg := range msgs {
				logger.Info(ctx, fmt.Sprintf("received message from %s", req.QueueName), zap.Any("data length", len(msg.Body)))

				err = req.Callback(ctx, msg.Body)
				if err != nil {
					fmt.Println("***************************")
					fmt.Println("Max Retry ", req.MaxRetries)
					err1 := r.Retry(ctx, msg, req.MaxRetries)
					if err1 != nil {
						errMsg := "failed from retry"
						logger.Error(ctx, fmt.Sprintf("%s : %s", errMsg, err1.Error()))
						return
					}
					errMsg := "failed from callback"
					logger.Error(ctx, fmt.Sprintf("%s : %s", errMsg, err.Error()))
					continue
				}
				msg.Ack(false)
			}
		}(ctx)

	}
}
