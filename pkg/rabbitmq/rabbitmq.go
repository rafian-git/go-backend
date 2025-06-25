package rabbitmq

import (
	"errors"
	ampq "github.com/rabbitmq/amqp091-go"
	"github.com/rafian-git/go-backend/pkg/log"
	"time"
)

type RabbitMQ struct {
	config  *Config
	log     *log.Logger
	retryIn int64 // ttl
}

func New(cnf *Config, logger *log.Logger) (Queue, error) {
	return &RabbitMQ{
		config: cnf,
		log:    logger.Named("rabbit-mq"),
	}, nil
}

func (r *RabbitMQ) Connect(uri string) (*ampq.Connection, error) {
	conn, err := ampq.Dial(uri)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (r *RabbitMQ) ConnectWithRetry(uri string, retries int, delay time.Duration) (*ampq.Connection, error) {
	var conn *ampq.Connection
	var err error
	for i := 0; i < retries; i++ {
		conn, err = r.Connect(uri)
		if err == nil {
			return conn, nil
		}
		time.Sleep(delay)
	}
	return nil, errors.New("failed to connect with RabbitMQ")
}

func (r *RabbitMQ) SetLogger(logger *log.Logger) {
	r.log = logger
}
