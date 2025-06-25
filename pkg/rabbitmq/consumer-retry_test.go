package rabbitmq

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.techetronventures.com/core/backend/pkg/log"
	"testing"
)

func porcessMessage(msg amqp.Delivery) bool {
	fmt.Println(msg)
	return false
}

func TestRabbitMQ_ConsumeWithRetry(t *testing.T) {
	type fields struct {
		config *Config
		log    *log.Logger
	}
	type args struct {
		ctx       context.Context
		ch        *amqp.Channel
		queueName string
	}
	conf := &Config{
		Url:     "amqp://guest:guest@localhost:5672",
		Retries: 0,
	}

	l := log.New().Named("test")

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    <-chan amqp.Delivery
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			fields: fields{
				config: conf,
				log:    l,
			},
			args: args{ctx: context.Background(), queueName: "test_qu"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RabbitMQ{
				config: tt.fields.config,
				log:    tt.fields.log,
			}
			conn, err := r.createConnection(context.Background())
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			ch, err := r.createChannel(context.Background(), conn)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			tt.args.ch = ch
			tt.args.queueName = "test_qu"
			fmt.Println("**********")

			got, err := r.ConsumeRMQ(tt.args.ctx, tt.args.ch, tt.args.queueName, 2000)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConsumeWithRetry() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for msg := range got {
				if !porcessMessage(msg) {
					r.Retry(tt.args.ctx, msg, 0)
				} else {
					msg.Ack(false)
				}
			}
			//
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("ConsumeWithRetry() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
