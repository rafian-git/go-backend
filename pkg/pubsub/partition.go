package pubsub

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func EnsureDesiredPartitions(ctx context.Context, config *Config, topic string, desiredPartitions int) error {

	conn, err := kafka.Dial("tcp", config.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return err
	}
	var partitionIds []int
	for _, partition := range partitions {
		if partition.Topic == topic {
			partitionIds = append(partitionIds, partition.ID)
		}
	}

	if len(partitionIds) == 0 {

		metadata, err := conn.Brokers()
		if err != nil {
			return err
		}

		fmt.Println(zap.Any("broker", metadata))
		err = conn.CreateTopics(kafka.TopicConfig{
			Topic:             topic,
			NumPartitions:     desiredPartitions,
			ReplicationFactor: len(metadata),
		})
		return err
	}

	if len(partitionIds) < desiredPartitions {
		client := kafka.Client{
			Addr: kafka.TCP(config.Address),
		}
		_, err := client.CreatePartitions(context.Background(), &kafka.CreatePartitionsRequest{
			Addr: kafka.TCP(config.Address),
			Topics: []kafka.TopicPartitionsConfig{
				{
					Name:  topic,
					Count: int32(desiredPartitions),
				},
			},
			ValidateOnly: false,
		})
		return err
	}

	return nil
}
