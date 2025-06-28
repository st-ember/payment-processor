package factory

import (
	"context"
	"paymentprocessor/internal/infra/config"
	kafkaadapter "paymentprocessor/internal/infra/kafka"
	"paymentprocessor/internal/infra/persistence/mongodb"
	redisadapter "paymentprocessor/internal/infra/redis"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Infrastructure holds all infrastructure components
type Infrastructure struct {
	MongoDB *mongo.Database
	Redis   *redis.Client
	Kafka   *kafkaadapter.KafkaClient
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
	kafkaClient, err := kafkaadapter.NewKafkaClient(cfg.Kafka.Brokers)
	if err != nil {
		return nil, err
	}

	return &Infrastructure{
		MongoDB: db,
		Redis:   redisClient,
		Kafka:   kafkaClient,
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
	return i.Kafka
}

func (i *Infrastructure) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := i.Redis.Close(); err != nil {
		return err
	}

	if err := i.Kafka.Close(); err != nil {
		return err
	}

	return i.MongoDB.Client().Disconnect(ctx)
}
