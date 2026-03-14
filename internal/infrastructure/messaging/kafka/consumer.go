package kafka

import (
	"context"
	"log"

	"github.com/twmb/franz-go/pkg/kgo"
)

type MessageHandler func(ctx context.Context, key string, value []byte) error

type Consumer struct {
	client *kgo.Client
}

func NewConsumer(brokers []string, groupID string, topics []string) (*Consumer, error) {
	cl, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topics...),
	)
	if err != nil {
		return nil, err
	}

	return &Consumer{client: cl}, nil
}

func (c *Consumer) Start(ctx context.Context, handler MessageHandler) {
	for {
		fetches := c.client.PollFetches(ctx)

		if errs := fetches.Errors(); len(errs) > 0 {
			for _, err := range errs {
				log.Println("Kafka error:", err)
			}
			continue
		}

		fetches.EachRecord(func(record *kgo.Record) {
			err := handler(ctx, string(record.Key), record.Value)
			if err != nil {
				log.Println("Handler error:", err)
			}
		})
	}
}

func (c *Consumer) Close() {
	c.client.Close()
}
