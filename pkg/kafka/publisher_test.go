package kafka

import (
	"fmt"
	km "github.com/segmentio/kafka-go"
	"gitlab.techetronventures.com/core/backend/pkg/log"
	"testing"
)

func TestPublishMessage(t *testing.T) {
	type args struct {
		conf    *Config
		message string
		topic   string
		n       int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test-1",
			args: args{
				conf: &Config{
					Address:                "localhost:9092",
					AllowAutoTopicCreation: false,
				},
				topic:   "test",
				n:       1,
				message: "Test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//writer := GetWriter(tt.args.conf)

			k := New(tt.args.conf, log.New())
			bytes, _ := km.Marshal(tt.args.message)
			//EnsureDesiredPartitions(context.Background(), tt.args.conf, tt.args.topic, 5)
			for i := 0; i < tt.args.n; i++ {
				err := k.Publish(tt.args.topic, km.Message{
					Topic: tt.args.topic,
					Value: bytes,
					Key:   bytes,
				})
				if err != nil {
					fmt.Errorf(err.Error())
				}
				//PublishMessage(context.Background(), writer, kafka.Message{Topic: tt.args.topic, Key: msg, Value: msg})
			}
			//writer.Close()
		})
	}
}
