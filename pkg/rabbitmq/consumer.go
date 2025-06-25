package rabbitmq

import (
	"context"
	"fmt"
	ampq "github.com/rabbitmq/amqp091-go"
)

func (r *RabbitMQ) Consume(ctx context.Context, ch *ampq.Channel, queueName string) (<-chan ampq.Delivery, error) {

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		msg := "failed to declare a queue"
		r.log.Error(ctx, fmt.Sprintf("%s : %s", msg, err.Error()))
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer // why empty?
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	return msgs, nil
}
