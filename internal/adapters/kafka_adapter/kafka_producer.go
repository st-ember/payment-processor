package kafkaadapter

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

func CreateProducer() (*kafka.Producer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaAdd,
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}
