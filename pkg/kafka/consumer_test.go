package kafka

import (
	"context"
	"fmt"
	"github.com/rafian-git/go-backend/pkg/log"
	"testing"
)

func CallbackTest(ctx context.Context, request []byte) error {
	fmt.Println(request)
	return fmt.Errorf("error")
}

func TestKafka_Listen(t *testing.T) {
	type fields struct {
		config *Config
		log    *log.Logger
	}
	type args struct {
		topic    string
		group    string
		callback ListenCallback
	}
	conf := &Config{
		Address: "localhost:9092", AllowAutoTopicCreation: true,
	}

	logger := log.New().Named("test-consumer")

	k := New(conf, logger)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "test",
			fields: fields{
				config: k.config,
				log:    logger,
			},
			args: args{
				topic:    "test",
				group:    "test-1-grp",
				callback: CallbackTest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Kafka{
				config: tt.fields.config,
				logger: tt.fields.log,
			}
			k.Subscribe(tt.args.topic, tt.args.group, tt.args.callback)
		})
	}
}
