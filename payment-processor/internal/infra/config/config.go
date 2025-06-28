package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for our application
type Config struct {
	MongoDB MongoDBConfig
	Redis   RedisConfig
	Kafka   KafkaConfig
	Server  ServerConfig
	Stripe  StripeConfig
	JWT     JwtConfig
}

// MongoDBConfig holds MongoDB connection settings
type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

// RedisConfig holds Redis connection settings
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// KafkaConfig holds Kafka connection settings
type KafkaConfig struct {
	Brokers []string
}

// ServerConfig holds HTTP server settings
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type StripeConfig struct {
	Secret     string
	SuccessURL string
}

type JwtConfig struct {
	Secret []byte
}

func LoadEnvConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, fmt.Errorf("invalid REDIS_DB value: %w", err)
	}

	cfg := &Config{
		MongoDB: MongoDBConfig{
			URI:      os.Getenv("MONGO_URI"),
			Database: os.Getenv("MONGO_DBNAME"),
			Timeout:  10 * time.Second,
		},
		Redis: RedisConfig{
			Addr:     os.Getenv("REDIS_ADD"),
			Password: os.Getenv("REDIS_PW"),
			DB:       redisDb,
		},
		Kafka: KafkaConfig{
			Brokers: []string{os.Getenv("KAFKA_BROKER")},
		},
		Server: ServerConfig{
			Port:         os.Getenv("SERVER_PORT"),
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		Stripe: StripeConfig{
			Secret:     os.Getenv("STRIPE_SECRET"),
			SuccessURL: os.Getenv("STRIPE_SUCCESS_URL"),
		},
		JWT: JwtConfig{
			Secret: []byte(os.Getenv("JWT_SECRET")),
		},
	}

	return cfg, nil
}

// NewDefaultConfig returns a new Config with default values
func NewDefaultConfig() *Config {
	return &Config{
		MongoDB: MongoDBConfig{
			URI:      "mongodb://localhost:27017",
			Database: "ecommerce_payment_ms",
			Timeout:  10 * time.Second,
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		Kafka: KafkaConfig{
			Brokers: []string{"localhost:9092"},
		},
		Server: ServerConfig{
			Port:         ":8080",
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}
