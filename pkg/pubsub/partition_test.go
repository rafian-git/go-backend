package pubsub

import (
	"context"
	"testing"
)

func TestEnsureDesiredPartitions(t *testing.T) {
	type args struct {
		ctx               context.Context
		config            *Config
		topic             string
		desiredPartitions int
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
				ctx:               context.Background(),
				config:            &Config{Address: "localhost:9092", AllowAutoTopicCreation: true},
				topic:             "new-topic",
				desiredPartitions: 14,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := EnsureDesiredPartitions(tt.args.ctx, tt.args.config, tt.args.topic, tt.args.desiredPartitions); (err != nil) != tt.wantErr {
				t.Errorf("EnsureDesiredPartitions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
