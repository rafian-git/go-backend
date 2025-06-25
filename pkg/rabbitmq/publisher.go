package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rafian-git/go-backend/utility"
	"go.uber.org/zap"
	"sync"
	"time"
)

func (r *RabbitMQ) createConnection(ctx context.Context) (*amqp.Connection, error) {
	conn, err := amqp.Dial(r.config.Url)
	if err != nil {
		msg := "failed to connect with rabbitmq"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}
	return conn, nil
}

func (r *RabbitMQ) createChannel(ctx context.Context, connection *amqp.Connection) (*amqp.Channel, error) {
	channel, err := connection.Channel()
	if err != nil {
		msg := "failed to connect with channel"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}

	if err = channel.Confirm(false); err != nil {
		r.log.Error(ctx, fmt.Sprintf("channel could not be put into confirm mode: %v", err))
		return nil, err
	}

	return channel, err
}

func (r *RabbitMQ) bindQueue(ctx context.Context, channel *amqp.Channel, queueName, routingKey, exchange string) error {
	err := channel.QueueBind(
		queueName,  // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,
		nil,
	)
	if err != nil {
		r.log.Error(ctx, fmt.Sprintf("failed to bind queue: %s", err.Error()))
	}
	return err
}

func (r *RabbitMQ) declareQueue(ctx context.Context, channel *amqp.Channel, queueName string, arg map[string]interface{}) error {

	_, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		arg,       // arguments
	)

	if err != nil {
		msg := "failed to declare queue"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
	}

	return err
}
func (r *RabbitMQ) publishToRabbitMQ(ctx context.Context, channel *amqp.Channel, routingKey, exchange string, obj []byte) error {

	err := channel.PublishWithContext(ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         obj,
		},
	)
	if err != nil {
		r.log.Error(ctx, fmt.Sprintf("failed to publish message: %v", err))
	}
	return err
}

func (r *RabbitMQ) publishWithRetry(ctx context.Context, channel *amqp.Channel, routingKey, exchange string, message interface{}) error {
	bytes := utility.CovertObjToBytes(message)

	//confirmations := channel.NotifyPublish(make(chan amqp.Confirmation, 1))
	for attempt := 1; attempt <= MaxRetriesPublish; attempt++ {
		err := r.publishToRabbitMQ(ctx, channel, routingKey, exchange, bytes)
		if err != nil {
			r.log.Info(ctx, fmt.Sprintf("failed to publish message %v to %s (Attempt %d): %v", message, exchange, attempt, err))
			time.Sleep(time.Second * time.Duration(attempt)) // Backoff: wait longer with each attempt
			continue
		}
		r.log.Info(ctx, fmt.Sprintf("published successfully to %s", routingKey), zap.Any("data-length", len(bytes)))
		return nil
		//select {
		//case confirm := <-confirmations:
		//	if confirm.Ack {
		//		r.log.Info(ctx, "Message published successfully to RabbitMQ", zap.Any("data-length", len(bytes)))
		//		return nil
		//	} else {
		//		r.log.Error(ctx, "Failed to publish message to RabbitMQ")
		//	}
		//case <-time.After(5 * time.Second):
		//	r.log.Error(ctx, fmt.Sprintf("failed to receive confirmation after %d attempts", attempt))
		//}
	}

	return fmt.Errorf("failed to publish message '%v' after %d attempts", zap.Any("msg", message), MaxRetriesPublish)
}

func (r *RabbitMQ) createExchange(ctx context.Context, channel *amqp.Channel, exchange, kind string) error {
	err := channel.ExchangeDeclare(
		exchange, // name
		kind,     // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		r.log.Error(ctx, fmt.Sprintf("failed to declare an exchange: %v", err))
	}
	return err
}

func (r *RabbitMQ) Publish(ctx context.Context, queueName string, obj []byte) error {
	// maybe we should not create a connection everytime we publish
	conn, err := amqp.Dial(r.config.Url)

	if err != nil {
		msg := "failed to connect with rabbitmq"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return err
	}

	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		msg := "failed to open a channel"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return err
	}

	defer ch.Close()

	args := make(amqp.Table)
	args["x-dead-letter-exchange"] = queueName + "-dlx"
	args["x-dead-letter-exchange"] = fmt.Sprintf("%s-dlx", queueName)

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments
	)

	if err != nil {
		msg := "failed to declare a queue"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return err
	}

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         obj,
		})

	if err != nil {
		msg := "failed to declare a queue"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return err
	}
	r.log.Info(ctx, "sent to rabbit-mq ", zap.Any("data-length ", len(obj)))
	return nil
}

func (r *RabbitMQ) PublishWithRetry(ctx context.Context, queueName string, obj []byte) error {
	conn, err := r.createConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	channel, err := r.createChannel(ctx, conn)
	if err != nil {
		return err
	}
	defer channel.Close()

	exchange := queueName + "-exchange"
	routingKey := queueName + "-routing-key"

	err = r.createExchange(ctx, channel, exchange, "direct")
	if err != nil {
		return err
	}

	delayExchange := queueName + "-dlx"
	err = r.declareQueue(ctx, channel, queueName, map[string]interface{}{
		"x-dead-letter-exchange": delayExchange,
	})

	if err != nil {
		return err
	}

	err = r.bindQueue(ctx, channel, queueName, routingKey, exchange)
	if err != nil {
		return err
	}

	for attempt := 1; attempt <= MaxRetriesPublish; attempt++ {
		err := r.publishToRabbitMQ(ctx, channel, routingKey, exchange, obj)
		if err != nil {
			r.log.Info(ctx, fmt.Sprintf("failed to publish message byte-len %d to %s (Attempt %d): %v", len(obj), exchange, attempt, err))
			time.Sleep(time.Second * time.Duration(attempt)) // Backoff: wait longer with each attempt
			continue
		}
		r.log.Info(ctx, fmt.Sprintf("published successfully to %s", routingKey), zap.Any("data-length", len(obj)))
		return nil
	}

	r.log.Info(ctx, "sent to rabbit-mq ", zap.Any("data-length ", len(obj)))
	return nil
}

func (r *RabbitMQ) PublishBulkMessage(ctx context.Context, queueName string, messages []interface{}) error {

	conn, err := r.createConnection(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	channel, err := r.createChannel(ctx, conn)
	if err != nil {
		return err
	}
	defer channel.Close()

	exchange := queueName + "-exchange"
	routingKey := queueName + "-routing-key"

	err = r.createExchange(ctx, channel, exchange, "direct")
	if err != nil {
		return err
	}

	delayExchange := queueName + "-dlx"
	err = r.declareQueue(ctx, channel, queueName, map[string]interface{}{
		"x-dead-letter-exchange": delayExchange,
	})

	if err != nil {
		return err
	}

	err = r.bindQueue(ctx, channel, queueName, routingKey, exchange)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, message := range messages {
		wg.Add(1)
		//TODO: Maybe need to add go routine
		func(msg interface{}) {
			defer wg.Done()
			if er := r.publishWithRetry(ctx, channel, routingKey, exchange, msg); er != nil {
				r.log.Error(ctx, fmt.Sprintf("err : %v", er))
			}
		}(message)
	}
	wg.Wait()
	return nil
}
