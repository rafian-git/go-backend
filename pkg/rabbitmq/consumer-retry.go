package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) ConsumeWithRetry(ctx context.Context, ch *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	delayQueue := queueName + "-dl-queue"
	exchange := queueName + "-exchange"
	routingKey := queueName + "-routing-key"
	delayExchange := queueName + "-dlx"

	///start dlx
	var restingQueue amqp.Queue
	restingQueue, err := ch.QueueDeclare(delayQueue, true, false, false, false, map[string]interface{}{
		"x-dead-letter-exchange":    exchange,
		"x-dead-letter-routing-key": routingKey,
		"x-max-length":              50000,
		"x-overflow":                "reject-publish",
		"x-message-ttl":             r.retryIn, //
	})

	if err != nil {
		msg := "failed to declare delay queue"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}
	err = ch.ExchangeDeclare(delayExchange, "topic", true, false, false, false, nil)
	if err != nil {
		msg := "failed to declare delay exchange"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}

	err = ch.QueueBind(restingQueue.Name, routingKey, delayExchange, false, nil)
	if err != nil {
		msg := "failed to bind dlx with dlq"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}
	///end dlx

	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		msg := "failed to declare exchange"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}

	var queue amqp.Queue
	queue, err = ch.QueueDeclare(queueName, false, false, false, false, map[string]interface{}{
		"x-dead-letter-exchange": delayExchange,
	})
	if err != nil {
		msg := "failed to declare queue"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}

	err = ch.QueueBind(queue.Name, routingKey, exchange, false, nil)
	if err != nil {
		msg := "failed to bind queue"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}

	// consumer name ???
	var deliveries <-chan amqp.Delivery
	deliveries, err = ch.Consume(queue.Name, "", false, false, false, false, nil)

	return deliveries, nil
}

func (r *RabbitMQ) ConsumeRMQ(ctx context.Context, ch *amqp.Channel, queueName string, ttl int64) (<-chan amqp.Delivery, error) {

	delayQueue := queueName + "-dl-queue"
	exchange := queueName + "-exchange"
	routingKey := queueName + "-routing-key"
	delayExchange := queueName + "-dlx"

	///start dlx
	var restingQueue amqp.Queue
	restingQueue, err := ch.QueueDeclare(delayQueue, true, false, false, false, map[string]interface{}{
		"x-dead-letter-exchange":    exchange,
		"x-dead-letter-routing-key": routingKey,
		"x-max-length":              50000,
		"x-overflow":                "reject-publish",
		"x-message-ttl":             ttl, //
	})

	if err != nil {
		msg := "failed to declare delay queue"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}
	err = ch.ExchangeDeclare(delayExchange, "topic", true, false, false, false, nil)
	if err != nil {
		msg := "failed to declare delay exchange"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}

	err = ch.QueueBind(restingQueue.Name, routingKey, delayExchange, false, nil)
	if err != nil {
		msg := "failed to bind dlx with dlq"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}
	///end dlx

	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		msg := "failed to declare exchange"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}

	var queue amqp.Queue
	queue, err = ch.QueueDeclare(queueName, false, false, false, false, map[string]interface{}{
		"x-dead-letter-exchange": delayExchange,
	})
	if err != nil {
		msg := "failed to declare queue"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}

	err = ch.QueueBind(queue.Name, routingKey, exchange, false, nil)
	if err != nil {
		msg := "failed to bind queue"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}

	// consumer name ???
	var deliveries <-chan amqp.Delivery
	deliveries, err = ch.Consume(queue.Name, "", false, false, false, false, nil)
	return deliveries, nil
}

func (r *RabbitMQ) Retry(ctx context.Context, msg amqp.Delivery, maxRetry int64) error {
	fmt.Println("Retrying...")
	if msg.Headers["x-death"] != nil {
		for _, death := range msg.Headers["x-death"].([]interface{}) {
			fmt.Println("Entering into x-death......")
			deathMap := death.(amqp.Table)
			if deathMap["reason"] == "expired" {
				fmt.Println("x-death expired")
				count, ok := deathMap["count"].(int64)
				if ok && count <= maxRetry {
					fmt.Printf("retry count %v\n", count)
					_ = msg.Nack(false, false)
				} else {
					fmt.Printf("maximum retries has been exceeded: %v\n", msg.MessageId)
				}
			}
		}
	} else {
		fmt.Println("Nacking---")
		err := msg.Nack(false, false)
		if err != nil {
			fmt.Println("error while nacking ", err.Error())
			return err
		}
	}

	return nil
}
