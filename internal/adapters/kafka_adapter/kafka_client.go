package kafkaadapter

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaClient struct{}

func NewKafkaClient() *KafkaClient { return &KafkaClient{} }

func (a *KafkaClient) SendMessage(p *kafka.Producer, topic string, key, value []byte) error {
	deliveryChan := make(chan kafka.Event)

	err := p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}, deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	close(deliveryChan)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}

type MessageSender interface {
	SendMessage(p *kafka.Producer, topic string, key, value []byte) error
}
