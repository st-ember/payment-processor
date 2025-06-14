package kafkaadapter

import (
	"encoding/json"
	"time"

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

func (c *KafkaClient) SendMessage(topic string, value map[string]interface{}) error {
	valBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(valBytes),
	}

	_, _, err = c.producer.SendMessage(msg)
	return err
}

func (c *KafkaClient) LogError(topic, description string, err error) error {
	value := map[string]interface{}{
		"time_stamp": time.Now(),
		description:  err,
	}
	return c.SendMessage(topic, value)
}

func (c *KafkaClient) Close() error {
	return c.producer.Close()
}
