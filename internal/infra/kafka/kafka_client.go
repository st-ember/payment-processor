package kafkaadapter

import (
	"encoding/json"

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

func (a *KafkaClient) SendMessage(topic string, key []byte, value map[string]string) error {
	deliveryChan := make(chan kafka.Event)

	valBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = a.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Key:            key,
		Value:          valBytes,
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

func (a *KafkaClient) LogError(topic string, key []byte, value map[string]string) error {
	deliveryChan := make(chan kafka.Event)

	valBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = a.p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Key:            key,
		Value:          valBytes,
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
