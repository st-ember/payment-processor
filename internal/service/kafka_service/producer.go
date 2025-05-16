package kafkaservice

import (
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

var kafkaAdd string

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	kafkaAdd = os.Getenv("KAFKA_ADD")
}

type KafkaService struct{}

func NewKafkaService() *KafkaService { return &KafkaService{} }

func (s *KafkaService) CreateProducer() (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaAdd,
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *KafkaService) SendMessage(p *kafka.Producer, topic string, key, value []byte) error {
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

type ProducerCreator interface {
	CreateProducer() (*kafka.Producer, error)
}

type MessageSender interface {
	SendMessage(p *kafka.Producer, topic string, key, value []byte) error
}
