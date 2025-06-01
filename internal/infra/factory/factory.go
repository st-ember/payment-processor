package factory

import (
	"context"
	"paymentprocessor/internal/infra/config"
	kafkaadapter "paymentprocessor/internal/infra/kafka_adapter"
	"paymentprocessor/internal/infra/persistence/mongodb"
	redisadapter "paymentprocessor/internal/infra/redis_adapter"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Infrastructure holds all infrastructure components
type Infrastructure struct {
	MongoDB *mongo.Database
	Redis   *redis.Client
	Kafka   *kafka.Producer
	Config  *config.Config
}

// NewInfrastructure creates and initializes all infrastructure components
func NewInfrastructure(cfg *config.Config) (*Infrastructure, error) {
	// Initialize MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), cfg.MongoDB.Timeout)
	defer cancel()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDB.URI))
	if err != nil {
		return nil, err
	}

	db := mongoClient.Database(cfg.MongoDB.Database)

	// Initialize Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	if err := redisClient.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	// Initialize Kafka
	kafkaConfig := &kafka.ConfigMap{"bootstrap.servcers": cfg.Kafka.Brokers}

	// producer, err := kafka.NewProducer(cfg.Kafka.Brokers, kafkaConfig)
	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		return nil, err
	}

	return &Infrastructure{
		MongoDB: db,
		Redis:   redisClient,
		Kafka:   producer,
		Config:  cfg,
	}, nil
}

func (i *Infrastructure) NewSessionRepository() *mongodb.SessionRepo {
	return mongodb.NewSessionRepository(i.MongoDB)
}

func (i *Infrastructure) NewRedisUtil() *redisadapter.RedisUtil {
	return redisadapter.NewRedisUtil(i.Redis)
}

func (i *Infrastructure) NewKafkaClient() *kafkaadapter.KafkaClient {
	return kafkaadapter.NewKafkaClient(i.Kafka)
}

func (i *Infrastructure) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := i.Redis.Close(); err != nil {
		return err
	}

	// no error check for confluent kafka shutdown
	i.Kafka.Flush(15 * 1000)
	i.Kafka.Close()

	return i.MongoDB.Client().Disconnect(ctx)
}
