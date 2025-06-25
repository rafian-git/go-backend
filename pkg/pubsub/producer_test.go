package pubsub

import (
	"context"
	"fmt"
	"github.com/rafian-git/go-backend/utility"
	"github.com/segmentio/kafka-go"
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
					Address:                "100.101.66.10:9093",
					AllowAutoTopicCreation: true,
				},
				topic:   "topic-partition",
				n:       10,
				message: "Test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := GetSingletonWriter(tt.args.conf)
			fmt.Println(writer.Addr)
			msg := utility.CovertObjToBytes(tt.args.message)
			//EnsureDesiredPartitions(context.Background(), tt.args.conf, tt.args.topic, 5)
			for i := 0; i < tt.args.n; i++ {
				PublishMessage(context.Background(), writer, kafka.Message{Topic: tt.args.topic, Key: msg, Value: msg})
			}
			writer.Close()
		})

	}
}
