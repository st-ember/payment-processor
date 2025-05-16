package main

import (
	"net/http"
	"os"
	"paymentprocessor/internal/handler"
	"paymentprocessor/internal/mongo"
	"paymentprocessor/internal/service"
	kafkaservice "paymentprocessor/internal/service/kafka_service"
	"paymentprocessor/internal/service/stripeservice"

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

	kafkaService := kafkaservice.NewKafkaService()
	producer, err := kafkaService.CreateProducer()
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	orderRepo := mongo.NewOrderRepository(mongoClient, "e-commerce", "order")
	sessionService := stripeservice.NewCheckoutSessionService()
	paymentService := service.NewPaymentSerivce(orderRepo, sessionService)
	paymentHandler := handler.NewPaymentHandler(paymentService)

	http.HandleFunc("/payment/start", paymentHandler.PaymentStart)
	http.HandleFunc("/payment/confirm", handler.PaymentConfirmation)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}
}
