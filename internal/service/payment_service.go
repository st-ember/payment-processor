package service

import (
	"context"
	kafkaadapter "paymentprocessor/internal/adapters/kafka_adapter"
	redisadapter "paymentprocessor/internal/adapters/redis_adapter"
	stripeadapter "paymentprocessor/internal/adapters/stripe_adapter"
	"paymentprocessor/internal/enums"
	"paymentprocessor/internal/storage"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService struct {
	orderStatusUpdater storage.OrderStatusUpdater
	sessionStarter     stripeadapter.SessionStarter
	messageSender      kafkaadapter.MessageSender
	stringRecordStorer redisadapter.StringRecordStorer
	stringRecordGetter redisadapter.StringRecordGetter
	producer           *kafka.Producer
	rdb                *redis.Client
}

func NewPaymentSerivce(
	orderStatusUpdater storage.OrderStatusUpdater,
	sessionStarter stripeadapter.SessionStarter,
	messageSender kafkaadapter.MessageSender,
	redisAdapter interface {
		redisadapter.StringRecordStorer
		redisadapter.StringRecordGetter
	},
	p *kafka.Producer, rdb *redis.Client) *PaymentService {
	return &PaymentService{
		orderStatusUpdater: orderStatusUpdater,
		sessionStarter:     sessionStarter,
		messageSender:      messageSender,
		stringRecordStorer: redisAdapter,
		stringRecordGetter: redisAdapter,
		producer:           p,
		rdb:                rdb,
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
	err = s.messageSender.SendMessage(s.producer, "stripe.checkout_session", []byte(sessionId))
	if err != nil {
		return "", err
	}
	// tell kafka about the expected payment and the checkout id we got from stripe

	return "", nil
}

func (s *PaymentService) ConfirmPayment(ctx context.Context, sessionId string, status enums.OrderStatus) error {
	// get orderId from redis
	orderId, err := s.stringRecordGetter.GetStringRecord(*s.rdb, ctx, sessionId)
	if err != nil {
		return err
	}
	// confirm orderId exists in db

	// set order's status to the right one

	return nil
}

type PaymentProcessor interface {
	ProcessPayment(ctx context.Context, orderId primitive.ObjectID, stripeParams []*stripe.CheckoutSessionLineItemParams) (string, error)
}

type PaymentConfirmer interface {
	ConfirmPayment(ctx context.Context, sessionId string, status enums.OrderStatus) error
}
