package main

import (
	"net/http"
	"os"
	"paymentprocessor/internal/handler"
	"paymentprocessor/internal/mongo"
	"paymentprocessor/internal/service"
	"paymentprocessor/internal/service/kafka"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	mongoClient, err := mongo.Connect(os.Getenv("MONGO_CONN_URI"))
	if err != nil {
		panic(err)
	}

	kafkaService := kafka.NewKafkaService()
	producer, err := kafkaService.CreateProducer()
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	orderRepo := mongo.NewOrderRepository(mongoClient, "e-commerce", "orders")
	productRepo := mongo.NewProductRepository(mongoClient, "e-commerce", "products")

	// sessionService := stripeservice.NewCheckoutSessionService()
	paymentService := service.NewPaymentSerivce(orderRepo, orderRepo, productRepo, producer)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	http.HandleFunc("/payment/start", paymentHandler.PaymentStart)
	http.HandleFunc("/payment/confirm", paymentHandler.PaymentConfirmation)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
