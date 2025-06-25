package rabbitmq

import (
	"context"
	"github.com/rafian-git/go-backend/pkg/log"
	"github.com/rafian-git/go-backend/utility"
	"testing"
)

func TestRabbitMQ_PublishBulkMessage(t *testing.T) {
	type fields struct {
		config *Config
		log    *log.Logger
	}
	type A struct {
		Name string
	}
	var list []interface{}

	list = append(list, &A{"adf"})
	//list = append(list, &A{"123"})

	l := log.New().Named("test")
	conf := &Config{
		Url:     "amqp://guest:guest@localhost:5672",
		Retries: 3,
	}
	type args struct {
		ctx       context.Context
		queueName string
		messages  []interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			fields: fields{
				log:    l,
				config: conf,
			},
			args: args{
				ctx:       context.Background(),
				queueName: "intra_day_order",
				messages:  list,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RabbitMQ{
				config: tt.fields.config,
				log:    tt.fields.log,
			}
			if err := r.PublishBulkMessage(tt.args.ctx, tt.args.queueName, tt.args.messages); (err != nil) != tt.wantErr {
				t.Errorf("PublishBulkMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRabbitMQ_PublishWithRetryMessage(t *testing.T) {
	type fields struct {
		config *Config
		log    *log.Logger
	}
	type A struct {
		Name string
	}
	var list []interface{}

	list = append(list, &A{"adf"})
	list = append(list, &A{"123"})

	l := log.New().Named("test")
	conf := &Config{
		Url:     "amqp://guest:guest@localhost:5672",
		Retries: 3,
	}
	type args struct {
		ctx       context.Context
		queueName string
		messages  []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			fields: fields{
				log:    l,
				config: conf,
			},
			args: args{
				ctx:       context.Background(),
				queueName: "test_qu",
				messages:  utility.CovertObjToBytes(&A{"adf"}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RabbitMQ{
				config: tt.fields.config,
				log:    tt.fields.log,
			}
			if err := r.PublishWithRetry(tt.args.ctx, tt.args.queueName, tt.args.messages); (err != nil) != tt.wantErr {
				t.Errorf("PublishBulkMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
