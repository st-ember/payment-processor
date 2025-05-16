package service

import (
	"context"
	"paymentprocessor/internal/enums"
	"paymentprocessor/internal/mongo"
	kafkaservice "paymentprocessor/internal/service/kafka_service"
	"paymentprocessor/internal/service/stripeservice"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stripe/stripe-go"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService struct {
	orderStatusUpdater mongo.OrderStatusUpdater
	sessionStarter     stripeservice.SessionStarter
	messageSender      kafkaservice.MessageSender
	producer           *kafka.Producer
}

func NewPaymentSerivce(
	orderStatusUpdater mongo.OrderStatusUpdater,
	sessionStarter stripeservice.SessionStarter,
	messageSender kafkaservice.MessageSender,
	p *kafka.Producer) *PaymentService {
	return &PaymentService{
		orderStatusUpdater: orderStatusUpdater,
		sessionStarter:     sessionStarter,
		messageSender:      messageSender,
		producer:           p,
	}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, orderId primitive.ObjectID, stripeParams []*stripe.CheckoutSessionLineItemParams) (string, error) {
	// set order to payment pending
	err := s.orderStatusUpdater.UpdateOrderStatus(ctx, orderId, enums.PaymentPending)
	if err != nil {
		return "", err
	}

	// tell stripe about the expected payment
	sessionId, err := s.sessionStarter.StartSession(stripeParams)
	if err != nil {
		return "", nil
	}
	err := s.messageSender.SendMessage(s.producer, "stripe.checkout_session", []byte(sessionId))
	if err != nil {
		return "", err
	}
	// tell kafka about the expected payment and the checkout id we got from stripe

	return "", nil
}

type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, orderId primitive.ObjectID, stripeParams []*stripe.CheckoutSessionLineItemParams) (string, error)
}
