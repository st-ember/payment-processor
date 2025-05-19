package service

import (
	"context"
	"encoding/json"
	kafkaadapter "paymentprocessor/internal/adapters/kafka_adapter"
	redisadapter "paymentprocessor/internal/adapters/redis_adapter"
	stripeadapter "paymentprocessor/internal/adapters/stripe_adapter"
	"paymentprocessor/internal/enums"
	"paymentprocessor/internal/storage"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentService struct {
	orderStatusUpdater storage.OrderStatusUpdater
	orderGetter        storage.OrderGetter
	productGetter      storage.ProductGetter
	sessionStarter     stripeadapter.SessionStarter
	messageSender      kafkaadapter.MessageSender
	stringRecordStorer redisadapter.StringRecordStorer
	stringRecordGetter redisadapter.StringRecordGetter
	producer           *kafka.Producer
	rdb                *redis.Client
}

func NewPaymentSerivce(
	orderStatusUpdater storage.OrderStatusUpdater,
	orderGetter storage.OrderGetter,
	productGetter storage.ProductGetter,
	sessionStarter stripeadapter.SessionStarter,
	messageSender kafkaadapter.MessageSender,
	redisAdapter interface {
		redisadapter.StringRecordStorer
		redisadapter.StringRecordGetter
	},
	p *kafka.Producer, rdb *redis.Client) *PaymentService {
	return &PaymentService{
		orderStatusUpdater: orderStatusUpdater,
		orderGetter:        orderGetter,
		productGetter:      productGetter,
		sessionStarter:     sessionStarter,
		messageSender:      messageSender,
		stringRecordStorer: redisAdapter,
		stringRecordGetter: redisAdapter,
		producer:           p,
		rdb:                rdb,
	}
}

func (s *PaymentService) ProcessPayment(
	ctx context.Context,
	orderId primitive.ObjectID,
	userId string) (string, error) {

	order, err := s.orderGetter.GetOrder(ctx, orderId)
	if err != nil {
		return "", err
	}

	productObjectId, err := primitive.ObjectIDFromHex(order.Product.String())
	if err != nil {
		return "", err
	}

	product, err := s.productGetter.GetProduct(ctx, productObjectId)
	if err != nil {
		return "", err
	}

	// construce stripe param for checkout request
	stripeParams := []*stripe.CheckoutSessionLineItemParams{
		{
			Price:    stripe.String(product.StripePriceId),
			Quantity: stripe.Int64(1),
		},
	}

	// tell stripe about the expected payment
	checkoutId, checkoutUrl, err := s.sessionStarter.StartSession(stripeParams)
	if err != nil {
		return "", err
	}

	// set order to payment pending
	err = s.orderStatusUpdater.UpdateOrderStatus(
		ctx, orderId,
		enums.PaymentPending)
	if err != nil {
		return "", err
	}

	// tell kafka about the checkout id we got from stripe
	// checkoutId as key, orderId as value
	// for cleaner abstraction to check after we get webhook reply
	err = s.messageSender.SendMessage(
		s.producer,
		"stripe.checkout_session",
		[]byte(checkoutId),
		[]byte(orderId.Hex()))
	if err != nil {
		return "", err
	}

	logVal := map[string]string{
		"level":     "info",
		"timestamp": time.Now().Format(time.RFC3339),
		"message":   "Stripe session started",
		"productId": product.Id.String(),
		"userId":    userId,
	}

	logValBytes, err := json.Marshal(logVal)
	if err != nil {
		return "", err
	}

	// log stripe session event to kafka to store in db
	err = s.messageSender.SendMessage(
		s.producer,
		"logs.payment.checkout",
		[]byte(orderId.Hex()),
		logValBytes,
	)
	return checkoutUrl, nil
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
	ProcessPayment(
		ctx context.Context,
		orderId primitive.ObjectID,
		userId string) (string, error)
}

type PaymentConfirmer interface {
	ConfirmPayment(ctx context.Context, sessionId string, status enums.OrderStatus) error
}
