package kafkaadapter

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaClient struct {
	p *kafka.Producer
}

func NewKafkaClient(p *kafka.Producer) *KafkaClient {
	return &KafkaClient{
		p: p,
	}
}

func (a *KafkaClient) SendMessage(topic string, key, value []byte) error {
	deliveryChan := make(chan kafka.Event)

	err := a.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
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
	SendMessage(topic string, key, value []byte) error
}
