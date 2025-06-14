package kafkaadapter

import (
	"encoding/json"

	"github.com/IBM/sarama"
)

type KafkaClient struct {
	producer sarama.SyncProducer
}

func NewKafkaClient(brokers []string) (*KafkaClient, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaClient{
		producer: producer,
	}, nil
}

func (a *KafkaClient) SendMessage(topic string, value map[string]interface{}) error {
	valBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(valBytes),
	}

	_, _, err = a.producer.SendMessage(msg)
	return err
}

func (a *KafkaClient) LogError(topic string, value map[string]interface{}) error {
	return a.SendMessage(topic, value)
}

func (a *KafkaClient) Close() error {
	return a.producer.Close()
}
