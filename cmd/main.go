package main

import (
	"log"
	"net/http"
	httpDelivery "paymentprocessor/internal/delivery/http"
	"paymentprocessor/internal/infra/config"
	"paymentprocessor/internal/infra/factory"
	"paymentprocessor/internal/lib/jwt"
	stripeadapter "paymentprocessor/internal/lib/stripe"
	"paymentprocessor/internal/usecase"
	"paymentprocessor/internal/worker"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	cfg, err := config.LoadEnvConfig()
	if err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	infra, err := factory.NewInfrastructure(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize infrastructure: %v", err)
	}
	defer infra.Close()
	// init infra components
	sessionRepo := infra.NewSessionRepository()
	redisUtil := infra.NewRedisUtil()
	kafkaClient := infra.NewKafkaClient()

	// lib
	checkoutSessionUtil := stripeadapter.NewCheckoutSessionUtil()
	jwtHelper := jwt.NewJWTHelper(cfg.JWT.Secret)

	// usecase
	paymentUsecase := usecase.NewPaymentUsecase(sessionRepo, redisUtil, kafkaClient, checkoutSessionUtil)

	// handlers
	paymentHandler := httpDelivery.NewPaymentHandler(paymentUsecase, jwtHelper)

	http.HandleFunc("/payment/start", paymentHandler.PaymentStart)
	http.HandleFunc("/payment/confirm", paymentHandler.PaymentConfirmation)

	if err := http.ListenAndServe(":8000", nil); err != nil {
		panic(err)
	}

	// worker
	scheduler := worker.NewScheduler(sessionRepo)
	scheduler.Start()
}
