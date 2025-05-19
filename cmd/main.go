package main

import (
	"context"
	"net/http"
	"os"
	kafkaadapter "paymentprocessor/internal/adapters/kafka_adapter"
	redisadapter "paymentprocessor/internal/adapters/redis_adapter"
	stripeadapter "paymentprocessor/internal/adapters/stripe_adapter"
	"paymentprocessor/internal/handler"
	"paymentprocessor/internal/service"
	"paymentprocessor/internal/storage"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	mongoClient, err := storage.Connect(os.Getenv("MONGO_CONN_URI"))
	if err != nil {
		panic(err)
	}

	// kafka
	producer, err := kafkaadapter.CreateProducer()
	if err != nil {
		panic(err)
	}
	defer producer.Close()
	kafkaClient := kafkaadapter.NewKafkaClient()

	// redis
	redisClientFactory := redisadapter.NewRedisClientFactory()
	rdb, err := redisClientFactory.InitRedisClient(ctx)
	if err != nil {
		panic(err)
	}
	redisUtil := redisadapter.NewRedisUtil()

	// repo
	orderRepo := storage.NewOrderRepository(mongoClient, "e-commerce", "order")
	productRepo := storage.NewProductRepository(mongoClient, "e-commerce", "product")

	sessionService := stripeadapter.NewCheckoutSessionService()
	paymentService := service.NewPaymentSerivce(orderRepo, orderRepo, productRepo, sessionService, kafkaClient, redisUtil, producer, rdb)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	http.HandleFunc("/payment/start", paymentHandler.PaymentStart)
	http.HandleFunc("/payment/confirm", paymentHandler.PaymentConfirmation)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
