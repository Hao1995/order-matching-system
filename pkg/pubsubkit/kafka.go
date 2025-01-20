package pubsubkit

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

// KafkaPubSub implements both Publisher and Subscriber
type KafkaPubSub struct {
	// @todo: allow to use a custom logger

	Writer *kafka.Writer
	Reader *kafka.Reader
}

// NewKafkaPublisher creates a new Kafka publisher
func NewKafkaPublisher(brokers []string) Publisher {
	return &KafkaPubSub{
		Writer: &kafka.Writer{
			Addr: kafka.TCP(brokers...),
		},
	}
}

// Publish sends a message to a Kafka topic
func (k *KafkaPubSub) Publish(topic string, key, value []byte) error {
	return k.Writer.WriteMessages(context.Background(), kafka.Message{
		Topic: topic,
		Key:   key,
		Value: value,
	})
}

// NewKafkaSubscriber creates a new Kafka subscriber
func NewKafkaSubscriber(brokers []string, topic string, groupID string) Subscriber {
	return &KafkaPubSub{
		Reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			GroupID: groupID,
			Topic:   topic,
		}),
	}
}

// Subscribe listens to messages from a Kafka topic
func (k *KafkaPubSub) Subscribe(topic string, handler func(key, value []byte) error) error {
	k.Reader.SetOffset(kafka.LastOffset)

	for {
		message, err := k.Reader.FetchMessage(context.Background())
		if err != nil {
			return err
		}
		if err := handler(message.Key, message.Value); err != nil {
			log.Printf("Error processing message: %v", err)
		}

		if err := k.Reader.CommitMessages(context.Background(), message); err != nil {
			log.Printf("Error committing message: %v", err)
		}
	}
}

// Close closes Kafka reader and writer
func (k *KafkaPubSub) Close() error {
	if k.Reader != nil {
		if err := k.Reader.Close(); err != nil {
			return err
		}
	}

	if k.Writer != nil {
		if err := k.Writer.Close(); err != nil {
			return err
		}
	}

	return nil
}
